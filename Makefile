.PHONY: dev dev-frontend dev-backend dev-widgets test test-frontend test-backend test-widgets build clean

dev: ## Start all services via Docker Compose
	docker compose up --build

dev-frontend: ## Start frontend dev server
	cd frontend && npm run dev

dev-backend: ## Start backend dev server
	cd backend && go run ./cmd/server

dev-widgets: ## Start widget service
	cd widgets && uvicorn app.main:app --reload

test: test-frontend test-backend test-widgets ## Run all tests

test-frontend: ## Run frontend tests
	cd frontend && npm test

test-backend: ## Run backend tests
	cd backend && go test ./...

test-widgets: ## Run widget tests
	cd widgets && python -m pytest tests/

build: ## Build all for production
	cd frontend && npm run build
	cd backend && CGO_ENABLED=0 go build -o bin/server ./cmd/server

clean: ## Clean build artifacts
	rm -rf frontend/dist backend/bin

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
