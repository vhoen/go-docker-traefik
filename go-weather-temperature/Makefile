DOCKER_COMPOSE_BIN = docker-compose
DOCKER_COMPOSE = $(DOCKER_COMPOSE_BIN)

GO_BIN = go
GO = $(GO_BIN)

APP = docker exec -ti $$(docker-compose ps -q app)

all: start

start: stop-traefik start-traefik stop-logs start-logs
	@$(DOCKER_COMPOSE) up -d app

stop:
	@$(DOCKER_COMPOSE) stop app || true

start-traefik:
	@$(DOCKER_COMPOSE) up -d traefik

restart: stop start

reload: stop
	@$(DOCKER_COMPOSE) up --no-recreate -d app

# stop traefik if running and rm container
stop-traefik:
	@docker stop traefik || true && docker rm traefik || true

start-logs:
	@$(DOCKER_COMPOSE) up -d logs

stop-logs:
	@docker stop logs || true && docker rm logs || true

# Bash in app container
app-bash-console:
	$(APP) bash