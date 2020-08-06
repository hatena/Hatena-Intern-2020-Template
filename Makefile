export DOCKER_BUILDKIT=1

.PHONY: up
up:
	skaffold dev --cleanup=false
