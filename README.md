# 🔌 Figest-IntegracaoService

> ⚠️ **Educational Project Notice**: This service is part of the **Figest** financial ecosystem, created for study, research, and testing purposes to demonstrate Open Finance integrations and OFX file parsing in Go.

---

## 📌 Overview

**Figest-IntegracaoService** is a Go microservice built with Fiber. It handles banking statement imports (`.OFX` file parsing) and Open Finance API integrations (Pluggy Connect token issuance and webhooks).

---

## 🛠️ Tech Stack
* **Language:** Go 1.22
* **Web Framework:** Fiber v2
* **ORM:** GORM
* **Messaging / Queue:** RabbitMQ (driver `amqp091-go`)

---

## 🔌 API Endpoints

| Resource | Method | Endpoint | Description |
|---|---|---|---|
| **OFX Import** | `POST` | `/import/ofx` | Parse uploaded `.OFX` bank statement file |
| **Open Finance** | `GET` | `/pluggy/token` | Get Pluggy Connect token |
| | `GET` | `/pluggy/accounts` | Fetch connected bank accounts |
| | `POST` | `/pluggy/webhook` | Process Pluggy transaction webhooks |

---

## 🚀 Running Locally

```bash
go mod download
go run ./cmd/server/main.go
```
