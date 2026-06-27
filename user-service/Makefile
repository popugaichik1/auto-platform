include .env
export


export PROJECT_ROOT=$(shell pwd)


.DEFAULT_GOAL := help


env-up: ## env: Запустить окружение проекта
	@docker compose up -d user-postgres


env-down: ## env: Остановить окружение проекта
	@docker compose down user-postgres

kafka-up:
	@docker compose up -d kafka1 kafka2 kafka3

env-cleanup: ## env: Очистить окружение проекта
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down user-postgres port-forwarder && \
		rm -rf ./out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

logs-cleanup: ## env: Очистить файлы логов из out/logs
	@read -p "Очистить все log файлы? Опасность утери логов. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ./out/logs && \
		echo "Файлы логов очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi

env-port-forward: ## env: Открыть порты сервисов окружения
	@docker compose up -d port-forwarder

env-port-close: ## env: Закрыть порты сервисов окружения
	@docker compose down port-forwarder


run-server:
	@echo "Starting server..."
	@go run cmd/api/main.go

run-server-deploy:
	@docker compose up -d --build user-service

down-server-deploy:
	@docker compose down user-service

migrate-up: ## PostgreSQL: Накатить миграции
	@make migrate-action action=up

migrate-down: ## PostgreSQL: Откатить миграции
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутсвует необходимый параметр action. Пример: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm user-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@user-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

show-logs:
	@docker compose logs -f