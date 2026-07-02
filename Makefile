COMPOSE_FILE := infra/compose/docker-compose.yml

.PHONY: dev dev-down dev-logs dev-reset test

dev:
	docker compose -f $(COMPOSE_FILE) up --build

dev-down:
	docker compose -f $(COMPOSE_FILE) down

dev-logs:
	docker compose -f $(COMPOSE_FILE) logs -f

dev-reset:
	docker compose -f $(COMPOSE_FILE) down -v

test:
	go -C services/backend test ./...
	pnpm --filter @essk/web lint
	pnpm --filter @essk/web typecheck
	pnpm --filter @essk/web build
	docker compose -f $(COMPOSE_FILE) config
