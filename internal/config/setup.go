package config

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"consumo-real-server/internal/application/bomba"
	combustivelapp "consumo-real-server/internal/application/combustivel"
	"consumo-real-server/internal/application/empresa"
	"consumo-real-server/internal/application/fornecedor"
	"consumo-real-server/internal/application/frentista"
	"consumo-real-server/internal/application/local"
	"consumo-real-server/internal/application/patrimonio"
	"consumo-real-server/internal/application/reservatorio"
	"consumo-real-server/internal/application/seeds"
	"consumo-real-server/internal/application/unidadeadministrativa"
	"consumo-real-server/internal/application/usuario"
	"consumo-real-server/internal/infrastructure/database"
	"consumo-real-server/internal/infrastructure/security"
	"consumo-real-server/internal/routes"
	applogger "consumo-real-server/internal/shared/logger"
)

const (
	dbMaxAttempts   = 10
	dbRetryBase     = 2 * time.Second
	shutdownTimeout = 5 * time.Second
)

type Config struct {
	ServerPort string
	DB         database.Config

	AdminBaseNome  string
	AdminBaseEmail string
	AdminBaseSenha string

	JWTSecret string
	TokenTTL  time.Duration
}

// Run orquestra a inicialização completa do sistema:
// config -> banco (com retry) -> migrações -> seeds -> servidor HTTP -> shutdown gracioso.
func Run() error {
	log := applogger.New()
	log.Info("iniciando inicializacao do sistema", "fase", "boot")

	if err := godotenv.Load(); err == nil {
		log.Info("arquivo .env carregado", "fase", "config")
	}

	cfg := loadConfig()
	log.Info("configuracao carregada",
		"fase", "config",
		"server_port", cfg.ServerPort,
		"db_host", cfg.DB.Host,
		"db_port", cfg.DB.Port,
		"db_name", cfg.DB.Name)

	db, err := connectWithRetry(log, cfg.DB)
	if err != nil {
		return fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}
	log.Info("banco de dados conectado",
		"fase", "database",
		"host", cfg.DB.Host,
		"port", cfg.DB.Port,
		"name", cfg.DB.Name)

	if err := database.AutoMigrate(db); err != nil {
		return fmt.Errorf("falha ao executar migracoes: %w", err)
	}
	log.Info("migracoes aplicadas com sucesso", "fase", "database")

	if err := seedAdminBase(db, cfg, log); err != nil {
		return fmt.Errorf("falha ao executar os seeds: %w", err)
	}
	log.Info("seeds executados com sucesso", "fase", "seed")

	combustivelRepo := database.NewCombustivelGORMRepository(db)
	combustivelService := combustivelapp.NewService(combustivelRepo)
	combustivelHandler := routes.NewCombustivelHandler(combustivelService)

	hasher := security.NewBcryptHasher()
	tokens := security.NewJWTManager(cfg.JWTSecret, cfg.TokenTTL)

	usuarioRepo := database.NewUsuarioGORMRepository(db)
	usuarioService := usuario.NewService(usuarioRepo, hasher, tokens)

	empresaRepo := database.NewEmpresaGORMRepository(db)
	empresaService := empresa.NewService(empresaRepo)

	unidadeRepo := database.NewUnidadeAdministrativaGORMRepository(db)
	unidadeService := unidadeadministrativa.NewService(unidadeRepo)

	localRepo := database.NewLocalGORMRepository(db)
	localService := local.NewService(localRepo)

	patrimonioRepo := database.NewPatrimonioGORMRepository(db)
	patrimonioService := patrimonio.NewService(patrimonioRepo)

	reservatorioRepo := database.NewReservatorioGORMRepository(db)
	reservatorioService := reservatorio.NewService(reservatorioRepo)

	bombaRepo := database.NewBombaGORMRepository(db)
	bombaService := bomba.NewService(bombaRepo)

	frentistaRepo := database.NewFrentistaGORMRepository(db)
	frentistaService := frentista.NewService(frentistaRepo)

	fornecedorRepo := database.NewFornecedorGORMRepository(db)
	fornecedorService := fornecedor.NewService(fornecedorRepo)

	r := routes.NewRoutes(routes.Handlers{
		Combustivel:  combustivelHandler,
		Usuario:      routes.NewUsuarioHandler(usuarioService),
		Auth:         routes.NewAuthHandler(usuarioService),
		Empresa:      routes.NewEmpresaHandler(empresaService),
		UnidadeAdmin: routes.NewUnidadeAdministrativaHandler(unidadeService),
		Local:        routes.NewLocalHandler(localService),
		Patrimonio:   routes.NewPatrimonioHandler(patrimonioService),
		Reservatorio: routes.NewReservatorioHandler(reservatorioService),
		Bomba:        routes.NewBombaHandler(bombaService),
		Frentista:    routes.NewFrentistaHandler(frentistaService),
		Fornecedor:   routes.NewFornecedorHandler(fornecedorService),
	}, tokens)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Info("iniciando servidor http", "fase", "http", "endereco", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		log.Info("sinal de encerramento recebido, iniciando shutdown gracioso", "fase", "shutdown")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		log.Info("servidor encerrado com sucesso", "fase", "shutdown")
		return nil
	case err := <-errCh:
		return err
	}
}

// connectWithRetry tenta conectar e validar o banco com backoff, resiliente ao
// inicio concorrente dos containers no Docker Compose.
func connectWithRetry(log *slog.Logger, cfg database.Config) (*gorm.DB, error) {
	var lastErr error
	for attempt := 1; attempt <= dbMaxAttempts; attempt++ {
		db, err := database.Open(cfg, log)
		if err == nil {
			sqlDB, dbErr := db.DB()
			if dbErr == nil {
				if pingErr := sqlDB.Ping(); pingErr == nil {
					return db, nil
				} else {
					err = pingErr
				}
			} else {
				err = dbErr
			}
		}
		lastErr = err
		log.Warn("banco de dados indisponivel, tentando novamente",
			"fase", "database",
			"tentativa", attempt,
			"max_tentativas", dbMaxAttempts,
			"proxima_tentativa_em", dbRetryBase,
			"erro", err)

		time.Sleep(dbRetryBase)
	}
	return nil, lastErr
}

func loadConfig() Config {
	return Config{
		ServerPort: envOr("SERVER_PORT", "8080"),
		DB: database.Config{
			Host:     envOr("DB_HOST", "localhost"),
			Port:     envOr("DB_PORT", "5432"),
			User:     envOr("DB_USER", "postgres"),
			Password: envOr("DB_PASSWORD", "postgres"),
			Name:     envOr("DB_NAME", "consumo_real"),
			SSLMode:  envOr("DB_SSLMODE", "disable"),
			Timezone: envOr("DB_TIMEZONE", "America/Sao_Paulo"),
		},
		AdminBaseNome:  envOr("ADMIN_BASE_NOME", "Administrador"),
		AdminBaseEmail: envOr("ADMIN_BASE_EMAIL", "admin@consumoreal.com.br"),
		AdminBaseSenha: envOr("ADMIN_BASE_SENHA", "admin123"),

		JWTSecret: envOr("JWT_SECRET", "consumo-real-server-dev-secret"),
		TokenTTL:  time.Duration(envIntOr("TOKEN_TTL_HORAS", 8)) * time.Hour,
	}
}

func envIntOr(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return fallback
}

func seedAdminBase(db *gorm.DB, cfg Config, log *slog.Logger) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminBaseSenha), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	repo := seeds.NewAdminBaseGORMRepository(db)
	if err := seeds.SeedAdminBase(repo, cfg.AdminBaseNome, cfg.AdminBaseEmail, string(hash)); err != nil {
		return err
	}
	log.Info("administrador base garantido", "fase", "seed", "email", cfg.AdminBaseEmail)
	return nil
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
