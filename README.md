# 💸 Internal Transfers API

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)

> A high-performance REST API for internal money transfers between accounts, built with **Go 1.25**, **PostgreSQL**, and **pgx/v5**.

---

## ✨ Features

- **Account Management** — Create accounts and query balances
- **Atomic Transfers** — Move funds between accounts with full transactional safety
- **Deadlock-Free** — Consistent lock ordering prevents database deadlocks
- **Input Validation** — Comprehensive request validation with meaningful error messages
- **Health Check** — Built-in `/health` endpoint for monitoring

---

## 🏗️ Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| Database | PostgreSQL 16 |
| DB Driver | pgx/v5 |
| Router | net/http (stdlib) |
| Migrations | Raw SQL |

---

## 📁 Project Structure

```
cmd/server/main.go              — Entry point & dependency wiring
internal/
  apperror/errors.go            — Domain error types
  config/config.go              — Env-based configuration
  database/postgres.go          — pgx/v5 connection pool
  database/txmanager.go         — Transaction manager
  dto/dto.go                    — Request/Response DTOs
  handler/                      — HTTP handlers, router, middleware
  model/model.go                — Domain models
  repository/                   — SQL data access layer
  service/                      — Business logic & interfaces
migrations/                     — SQL migration files
run_tests.sh                    — Bash integration test suite
```

---

## 🚀 Quick Start

### Prerequisites

| Tool | Install | Verify |
|---|---|---|
| **Go** 1.25+ | `brew install go` | `go version` |
| **PostgreSQL** 14+ | `brew install postgresql@16` | `psql --version` |
| **Make** | Pre-installed on macOS | `make --version` |

```bash
brew services start postgresql@16
```

### Setup & Run

```bash
# 1. Clone the repository
git clone https://github.com/<your-username>/internal-transfers-api.git
cd internal-transfers-api

# 2. Install Go dependencies
go mod download

# 3. Create database & apply migrations
make local-setup

# 4. Start the server
make run
```

Verify it's running:

```bash
curl http://localhost:8080/health
# → {"status":"ok"}
```

---

## ⚙️ Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | PostgreSQL username |
| `DB_PASSWORD` | `postgres` | PostgreSQL password |
| `DB_NAME` | `transaction_manager` | PostgreSQL database name |
| `SERVER_PORT` | `8080` | HTTP server port |

---

## 📡 API Reference

### Health Check

```
GET /health
```

**Response:** `200 OK`
```json
{ "status": "ok" }
```

---

### Create Account

```
POST /accounts
```

**Request Body:**
```json
{ "account_id": 1, "initial_balance": "1000.00" }
```

| Status | Meaning |
|---|---|
| `201` | Account created |
| `400` | Validation error |
| `409` | Account already exists |

---

### Get Account

```
GET /accounts/{account_id}
```

**Response:** `200 OK`
```json
{ "account_id": 1, "balance": "1000" }
```

| Status | Meaning |
|---|---|
| `400` | Invalid ID format |
| `404` | Account not found |

---

### Transfer Funds

```
POST /transactions
```

**Request Body:**
```json
{
  "source_account_id": 1,
  "destination_account_id": 2,
  "amount": "100.00"
}
```

| Status | Meaning |
|---|---|
| `201` | Transfer completed |
| `400` | Validation error |
| `404` | Account not found |
| `422` | Insufficient balance |

---

## 🧪 Running Tests

```bash
# Go unit tests (with race detector)
make test

# Integration test suite (requires running server)
# Terminal 1:
make local-db-drop && make local-setup
make run

# Terminal 2:
bash run_tests.sh
```

---

## 🛠️ Makefile Reference

| Command | Description |
|---|---|
| `make build` | Build binary to `bin/server` |
| `make run` | Build and run the server |
| `make test` | Run Go tests with race detector |
| `make lint` | Run staticcheck linter |
| `make clean` | Remove build artefacts |
| **Local DB** | |
| `make local-setup` | Full setup: create DB + apply migrations |
| `make local-db-create` | Create PostgreSQL database |
| `make local-db-drop` | Drop PostgreSQL database |
| `make local-migrate-up` | Apply migrations |

---

## 🔧 Troubleshooting

<details>
<summary><b>Click to expand common issues</b></summary>

| Error | Fix |
|---|---|
| `connection refused` | `brew services start postgresql@16` |
| `role "postgres" does not exist` | `psql -d postgres -c "CREATE USER postgres WITH SUPERUSER PASSWORD 'postgres';"` |
| `database "transaction_manager" does not exist` | `make local-db-create` |
| `Server is not reachable` (tests) | Start the server first in another terminal |
| `409 Conflict` on tests | Reset DB: `make local-db-drop && make local-setup`, then restart server |
| `relation "accounts" does not exist` | `make local-migrate-up` |

</details>
