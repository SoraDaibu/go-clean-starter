SERVICE := go-clean-starter
TEST_SERVICE := $(SERVICE)-test

DC := docker compose
DC_TEST := $(DC) -p $(TEST_SERVICE) -f docker-compose.test.yaml

# ─── Quick Start ───────────────────────────────────────────────────────────
quickstart:
	@if [ ! -f .env ]; then \
		cp .env.example .env && \
		echo "✅ Created .env from .env.example"; \
	else \
		echo "ℹ️  .env already exists, skipping copy"; \
	fi
	$(MAKE) up
	@echo ""
	@echo "🎉 Quickstart complete!"
	@echo "📍 API running at http://localhost:8080"
	@echo "💡 Run 'make test' to verify everything works"

# ─── docker compose lifecycle ─────────────────────────────────────────────
build:
	$(DC) build --no-cache

up:
	$(DC) up -d

log-api:
	$(DC) logs -f api

log-task-runner:
	$(DC) logs -f task-runner

down:
	$(DC) down

down-api:
	$(DC) down api

down-test:
	$(DC_TEST) down

# `|| true` ensures that the command doesn't fail if the image doesn't exist (so Makefile won't stop with an error).
docker-rmi:
	docker rmi go-clean-starter:latest || true
	docker rmi go-clean-starter-test-test || true

# ─── Testing with separate test compose ─────────────────────────────────────────
test:
	$(DC_TEST) run --rm --entrypoint "go test -v ./..." app
	$(MAKE) down-test

# ─── Database migrations ───────────────────────────────────────────────────────
migration:
	migrate create -ext sql -dir migration/sql -seq $(name)

migrate:
	migrate -path migration/sql -database "$(db)" up

psql:
	$(DC) exec postgres psql -U postgres -d go_clean_starter

db=postgres://postgres:example@localhost:5432/go_clean_starter?sslmode=disable
clean-migration:
	@if [ -z "$(version)" ]; then \
		echo "❌ Error: version argument is required."; \
		exit 1; \
	fi
	migrate -path migration/sql -database "$(db)" force $(version)

# ─── Auto generate ───────────────────────────────────────────────────────
sqlc:
	sqlc generate

oapi-codegen:
	oapi-codegen -generate models -package handler -o internal/http/handler/openapi_types.gen.go doc/api.yaml

# ─── Task ─────────────────────────────────────────────────────────────
import-items:
	$(DC) --profile task run --rm task-runner go run . task import --source-dir=$(or $(source-dir),./internal/task/item/data) $(if $(dry-run),--dry-run)

import-items-dry:
	$(DC) --profile task run --rm task-runner go run . task import --source-dir=$(or $(source-dir),./internal/task/item/data) --dry-run


# ─── Chore ─────────────────────────────────────────────────────────────
tree:
	tree --dirsfirst -I 'node_modules|vendor'
