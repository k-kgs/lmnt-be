.PHONY: up down migrate-up migrate-down seed sqlc build logs test

up:
	docker compose up -d db
	docker compose up --build -d api

down:
	docker compose down

migrate-up:
	docker compose run --rm migrate up

migrate-down:
	docker compose run --rm migrate down 1

seed:
	docker compose exec -T db psql -U kayam -d kayam < seed/seed.sql

sqlc:
	sqlc generate

build:
	go build ./...

logs:
	docker compose logs -f api

test:
	go test ./...
