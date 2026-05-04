# Simplified Stock Market Service

A reliable, high-availability stock market simulator built with **Go 1.26** and **PostgreSQL 18**.

## 🚀 Overview
This service simulates a simplified stock exchange where users can trade stocks through a central Bank. It is designed with a focus on **Data Integrity** (using ACID transactions) and **High Availability** (using a Load Balancer and multi-instance replication).

## 🏗️ System Architecture
The solution follows a modern, distributed approach:
- **Load Balancer (Nginx):** Acts as the single entry point, distributing traffic across multiple application instances.
- **Application Replicas (Go):** 3 independent instances of the service. If one fails (or is killed via `/chaos`), the system remains fully operational.
- **Database (PostgreSQL):** A single source of truth ensuring atomic operations for all trades.
- **CI/CD:** Automated testing and linting via GitHub Actions.

## 🛠️ Prerequisites
- **Docker** and **Docker Compose** installed.
- **Make** (for Linux/macOS/WSL) or **PowerShell** (for Windows).

## 🚦 Getting Started

### One-Command Startup
The application is fully containerized and works on all major operating systems (Windows, Linux, macOS) and architectures (x64, arm64).

**For Linux / macOS / WSL:**
```bash
make run XXXX=8888
```

**For Windows (PowerShell):**
```powershell
./run.ps1 -XXXX 8888
```
*The service will be available at `http://localhost:8888`.*

### Running Tests
To run the automated test suite without installing Go locally:
```bash
make test
```

---

## 📡 API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/stocks` | Sets the initial state of the Bank. |
| `GET` | `/stocks` | Returns the current state of the Bank. |
| `POST` | `/wallets/{id}/stocks/{name}` | Execute a trade (`buy` or `sell`). |
| `GET` | `/wallets/{id}` | Returns the current state of a specific wallet. |
| `GET` | `/wallets/{id}/stocks/{name}` | Returns the quantity of a specific stock in a wallet. |
| `GET` | `/log` | Returns the entire audit log (successful operations only). |
| `POST` | `/chaos` | Kills the instance serving the request to test HA. |

---

## 🧠 Design Decisions

- **PostgreSQL Transactions:** Every trade operation is wrapped in a database transaction. This ensures that if the Bank's stock is deducted, the Wallet's stock is increased atomically. No "lost" stocks.
- **Graceful Shutdown:** The service listens for termination signals and finishes processing active requests before shutting down, preventing data corruption.
- **Nginx & Replicas:** By using Nginx to load balance across 3 replicas, the system satisfies the High Availability requirement. Killing one instance via `/chaos` results in zero downtime for the user.
- **Interface-based Storage:** The Handlers depend on a `Storage` interface, not a concrete implementation. This allowed for **Table-Driven Testing** using Mocks.

## 📁 Project Structure
```text
├── cmd/server          # Application entry point & Graceful Shutdown logic
├── internal/
│   ├── handlers        # HTTP Handlers & Table-Driven Tests
│   ├── models          # Domain Data Structures (JSON tags)
│   └── storage         # PostgreSQL implementation & Transactions
├── nginx.conf          # Load Balancer configuration
├── docker-compose.yml  # Infrastructure orchestration (3 replicas)
├── Makefile            # Automation for Linux/WSL
└── run.ps1             # Automation for Windows
```
---
**Author:** Szymon Cimochowski 