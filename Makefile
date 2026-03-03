APP_NAME=accounts-api
DB_URL=postgres://postgres:postgres@localhost:5432/accounts?sslmode=disable

.PHONY: up down logs run test tidy migrate

up:
	docker compose up -d
	
down:
	docker compose down

logs:
	docker compose logs -f

tidy:
	go mod tidy

run:
	APP_DB_URL="$(DB_URL)" go run ./cmd/api

test:
	APP_DB_URL="$(DB_URL)" go test ./... -count=1

migrate:
	APP_DB_URL="$(DB_URL)" go run ./cmd/api -migrate