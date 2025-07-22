include scripts/*.mk

DEV_COMPOSE_ARGS=--env-file .env.dev -f docker-compose.dev.yaml
DEV_COMPOSE_ENV=docker compose $(DEV_COMPOSE_ARGS)
DEV_COMPOSE=docker compose $(DEV_COMPOSE_ARGS)


dev-local-build:
	$(DEV_COMPOSE) build
dev-local-up: dev-local-build
	$(DEV_COMPOSE) --env-file .env.local up -d
dev-build: dev-api-build
	$(DEV_COMPOSE) build
dev-up: dev-build dev-api-up
	$(DEV_COMPOSE) up -d
dev-api-build: api_docker_build

dev-api-up:
	$(DEV_COMPOSE) -f docker-compose.dev.yaml up -d api-up
env-check:
	docker compose --env-file .env.dev config