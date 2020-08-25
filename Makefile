export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1

.PHONY: up
up:
	skaffold dev --cleanup=false

.PHONY: dc
dc:
	docker-compose up

.PHONY: dcu
dcu:
	docker-compose build
	docker-compose up -d
