# consumo-real-server

Sistema de gerenciamento de abastecimento de frota (multitenant por empresa).

Stack: Go + GORM (ORM) + PostgreSQL + Gorilla Mux.

## Como rodar

O servidor carrega automaticamente o arquivo `.env` da raiz do projeto (se
existir). Crie-o a partir do exemplo:

```bash
cp .env.example .env
```

> **Importante:** a API se conecta ao Postgres usando `DB_HOST`, `DB_PORT` etc.
> do `.env`. Se o banco estiver num host/porta diferente, ajuste o `.env`.

### Opção A — Tudo com Docker (recomendado)

Sobe o Postgres + a API juntos:

```bash
docker compose up --build -d
```

O compose sobe dois serviços:

- `db`: PostgreSQL 16 (publicado no host na porta `DB_PORT`, volume `pgdata`)
- `app`: a API Go (publicada no host na porta `SERVER_PORT`)

No primeiro boot o servidor executa automaticamente o AutoMigrate (criação das
tabelas a partir dos domains) e o seed do administrador base.

Verificação:

```bash
curl http://localhost:8080/health   # use a porta do seu .env
```

Logs:

```bash
docker compose logs -f app
```

### Opção B — API localmente, banco no Docker

Suba apenas o Postgres e rode a API com Go:

```bash
docker compose up -d db
go run ./cmd/server
```

O `.env` precisa apontar para `DB_HOST=localhost` e `DB_PORT` igual à porta
publicada no host (padrão `5432`, ou a definida no `.env`).

### Portas em uso por outro projeto?

Se `8080`/`5432` já estiverem ocupadas na sua máquina (ex.: outro container
rodando), defina portas livres no `.env`:

```dotenv
SERVER_PORT=8081
DB_PORT=5433
```

Internamente a API continua escutando em `8080` e o app acessa o banco em
`db:5432` (rede interna do compose), independentemente da porta publicada.

## Configuração (variáveis de ambiente)

| Variável | Padrão | Descrição |
| --- | --- | --- |
| `SERVER_PORT` | `8080` | Porta do servidor HTTP |
| `DB_HOST` | `localhost` | Host do PostgreSQL |
| `DB_PORT` | `5432` | Porta do PostgreSQL |
| `DB_USER` | `postgres` | Usuário do banco |
| `DB_PASSWORD` | `postgres` | Senha do banco |
| `DB_NAME` | `consumo_real` | Nome do banco |
| `DB_SSLMODE` | `disable` | Modo SSL da conexão |
| `DB_TIMEZONE` | `America/Sao_Paulo` | Fuso horário da conexão |
| `ADMIN_BASE_NOME` | `Administrador` | Nome do admin base |
| `ADMIN_BASE_EMAIL` | `admin@consumoreal.com.br` | E-mail do admin base |
| `ADMIN_BASE_SENHA` | `admin123` | Senha do admin base (armazenada com bcrypt) |

## Criação de empresas (multitenant)

O sistema é multitenant: cada empresa é um tenant e todo usuário (exceto o
`ADMIN_BASE`) pertence a uma empresa. Isso cria um ciclo no onboarding: para
criar uma empresa seria preciso já existir um usuário, e usuário exige empresa.
Para resolver, `POST /api/empresas` suporta dois fluxos:

### 1) Apenas a empresa (qualquer usuário autenticado)

Cria somente o registro da empresa. Útil para administradores de um tenant já
existente abrirem outra empresa.

```bash
curl -X POST http://localhost:8080/api/empresas \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"nome":"Posto Teste","cnpj":"00.000.000/0001-00"}'
```

### 2) Empresa + primeiro administrador (somente `ADMIN_BASE`)

O fluxo recomendado para "onboarding" de um novo tenant: em **uma única
transação atômica** são criados a empresa e o usuário administrador inicial
(papel `ADMINISTRADOR`). Se qualquer campo for inválido, nada é persistido.

```bash
curl -X POST http://localhost:8080/api/empresas \
  -H "Authorization: Bearer <token-do-admin-base>" \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Posto Teste",
    "cnpj": "00.000.000/0001-00",
    "nome_administrador": "João da Silva",
    "email_administrador": "joao@posto.com",
    "senha_administrador": "senha123"
  }'
```

O usuário criado já consegue fazer login em `POST /api/auth/login` com o
e-mail/senha informados.

> **Notas**
> - Os três campos do administrador devem ser informados juntos; informar
>   apenas alguns retorna erro de validação.
> - Somente o papel `ADMIN_BASE` pode usar o fluxo 2; outro papel autenticado
>   recebe `401`.
> - O `ADMIN_BASE` é criado automaticamente no boot (ver variáveis acima).

## Documentação (Swagger)

A documentação interativa da API é servida pelas rotas abaixo (públicas, sem
autenticação):

| Método | Rota | Descrição |
| --- | --- | --- |
| GET | `/swagger` | Interface Swagger UI |
| GET | `/swagger/` | Interface Swagger UI (variante com barra) |
| GET | `/swagger/doc.json` | Spec Swagger 2.0 em JSON |
