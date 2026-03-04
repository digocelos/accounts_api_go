# Accounts API (Go + Postgres)

A small, production-style Go service showcasing:
- REST API design
- Idempotent create (by `document`)
- Optimistic locking (`version`)
- Postgres persistence
- Structured logging
- Basic health endpoint

## Requirements
- Go 1.22+
- Docker / Docker Compose

## Run locally

```bash
make up
make migrate
make run