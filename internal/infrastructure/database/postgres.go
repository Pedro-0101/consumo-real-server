package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"consumo-real-server/internal/domain/abastecimento"
	"consumo-real-server/internal/domain/bomba"
	"consumo-real-server/internal/domain/combustivel"
	"consumo-real-server/internal/domain/empresa"
	"consumo-real-server/internal/domain/entrada"
	"consumo-real-server/internal/domain/fornecedor"
	"consumo-real-server/internal/domain/frentista"
	"consumo-real-server/internal/domain/local"
	"consumo-real-server/internal/domain/medicao"
	"consumo-real-server/internal/domain/ordemabastecimento"
	"consumo-real-server/internal/domain/patrimonio"
	"consumo-real-server/internal/domain/preco"
	"consumo-real-server/internal/domain/reservatorio"
	"consumo-real-server/internal/domain/unidadeadministrativa"
	"consumo-real-server/internal/domain/usuario"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Timezone string
}

func Open(cfg Config, log *slog.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode, cfg.Timezone,
	)

	var gormLog gormlogger.Interface
	if log != nil {
		gormLog = &slogGormLogger{logger: log.With("componente", "gorm"), level: gormlogger.Warn}
	} else {
		gormLog = gormlogger.Default.LogMode(gormlogger.Warn)
	}

	return gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLog})
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&empresa.Empresa{},
		&unidadeadministrativa.UnidadeAdministrativa{},
		&local.Local{},
		&patrimonio.Patrimonio{},
		&usuario.Usuario{},
		&combustivel.Combustivel{},
		&reservatorio.Reservatorio{},
		&bomba.Bomba{},
		&bomba.Bico{},
		&frentista.Frentista{},
		&fornecedor.Fornecedor{},
		&preco.Preco{},
		&ordemabastecimento.OrdemAbastecimento{},
		&entrada.Entrada{},
		&medicao.Medicao{},
		&abastecimento.Abastecimento{},
	)
}

// slogGormLogger adapta o logger do GORM para os logs estruturados em JSON do serviço.
type slogGormLogger struct {
	logger *slog.Logger
	level  gormlogger.LogLevel
}

func (s *slogGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return &slogGormLogger{logger: s.logger, level: level}
}

func (s *slogGormLogger) Info(_ context.Context, msg string, data ...interface{}) {
	if s.level >= gormlogger.Info {
		s.logger.Info(msg, toAttrs(data)...)
	}
}

func (s *slogGormLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	if s.level >= gormlogger.Warn {
		s.logger.Warn(msg, toAttrs(data)...)
	}
}

func (s *slogGormLogger) Error(_ context.Context, msg string, data ...interface{}) {
	if s.level >= gormlogger.Error {
		s.logger.Error(msg, toAttrs(data)...)
	}
}

func (s *slogGormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if s.level == gormlogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	switch {
	case err != nil && s.level >= gormlogger.Error:
		s.logger.Error("erro no banco de dados",
			"sql", sql, "linhas", rows, "duracao", elapsed, "erro", err)
	case elapsed > time.Second && s.level >= gormlogger.Warn:
		s.logger.Warn("consulta lenta",
			"sql", sql, "linhas", rows, "duracao", elapsed)
	case s.level >= gormlogger.Info:
		s.logger.Info("consulta executada",
			"sql", sql, "linhas", rows, "duracao", elapsed)
	}
}

func toAttrs(data []interface{}) []any {
	attrs := make([]any, 0, len(data))
	for i := 0; i+1 < len(data); i += 2 {
		attrs = append(attrs, data[i], data[i+1])
	}
	return attrs
}
