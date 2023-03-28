# https://gist.github.com/tadashi-aikawa/da73d277a3c1ec6767ed48d1335900f3
.PHONY: $(shell egrep -oh ^[a-zA-Z0-9][a-zA-Z0-9\/_-]+: $(MAKEFILE_LIST) | sed 's/://')

all: upd ps

ps: ## 起動中のコンテナを表示
	docker compose ps
up: ## コンテナを起動
	docker compose up
upd: ## コンテナをバックグラウンドで起動
	docker compose up -d
stop: ## コンテナを停止
	docker compose stop
down: ## コンテナを破棄
	docker compose down --remove-orphans
downv: ## コンテナとボリューム、ネットワークを破棄
	docker compose down -v --remove-orphans
prune: ## 不要なDockerイメージを破棄
	docker system prune -f
restart: ## コンテナを再起動
	docker compose restart
login: ## Mysqlコンテナへログイン
	docker compose exec mysql bash

# https://postd.cc/auto-documented-makefile/
help: ## ヘルプ
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9][a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
