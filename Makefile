# https://gist.github.com/tadashi-aikawa/da73d277a3c1ec6767ed48d1335900f3
.PHONY: $(shell egrep -oh ^[a-zA-Z0-9][a-zA-Z0-9\/_-]+: $(MAKEFILE_LIST) | sed 's/://')

all: upd ps

ps: ## Display containers being started
	docker compose ps
up: ## Start containers
	docker compose up
upd: ## Start containers in background
	docker compose up -d
stop: ## Stop containers
	docker compose stop
down: ## Destroy containers
	docker compose down --remove-orphans
downv: ## Destroy containers, volumes and networks
	docker compose down -v --remove-orphans
prune: ## Destroy unneeded Docker images
	docker system prune -f
restart: ## Restart containers
	docker compose restart
login: ## Login to Mysql container
	docker compose exec mysql bash

# https://postd.cc/auto-documented-makefile/
help: ## Help
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9][a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
