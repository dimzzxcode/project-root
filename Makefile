.PHONY: help up down logs ps build restart migrate migrate-version dev-backend dev-frontend migrate-local install lint

help: ## Tampilkan daftar target
	@echo "Targets:"
	@echo "  up               Start semua service dengan Docker Compose (build ulang)"
	@echo "  down             Stop dan hapus container (data tetap)"
	@echo "  logs             Ikuti log semua service"
	@echo "  ps               Status container"
	@echo "  build            Build image tanpa start"
	@echo "  restart          Restart semua service"
	@echo "  migrate          Jalankan migration di container backend"
	@echo "  migrate-version  Tampilkan versi migration saat ini"
	@echo "  dev-backend      Jalankan backend secara lokal (butuh PostgreSQL)"
	@echo "  dev-frontend     Jalankan frontend secara lokal"
	@echo "  migrate-local    Jalankan migration secara lokal"
	@echo "  install          Install dependensi backend dan frontend"
	@echo "  lint             Jalankan go vet dan ESLint"

up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f

ps:
	docker compose ps

build:
	docker compose build

restart:
	docker compose restart

migrate:
	docker compose run --rm --build backend /app/migrate up

migrate-version:
	docker compose run --rm backend /app/migrate version

dev-backend:
	cd backend && go run ./cmd/api

dev-frontend:
	cd frontend && npm run dev

migrate-local:
	cd backend && go run ./cmd/migrate up

install:
	cd backend && go mod download
	cd frontend && npm install

lint:
	cd backend && go vet ./...
	cd frontend && npm run lint
