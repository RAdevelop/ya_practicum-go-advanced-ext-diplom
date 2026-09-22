# Переменные
## Это конечно не стоит хранить в Git.
export POSTGRES_USER := postgres
export POSTGRES_PASSWORD := postgres
export POSTGRES_HOST := postgres
export POSTGRES_PORT := 5432
export POSTGRES_DB := praktikum

export POSTGRES_DB_TEST := praktikum_test

export DB_DSN_TEST=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB_TEST}?sslmode=disable
DB_DSN=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable

# Определение команд
DOCKER_COMPOSE := docker-compose
GO_GENERATE := go generate ./internal/...
GO_TEST := go test ./internal/...

# Цвета для вывода
GREEN = \033[0;32m
YELLOW = \033[0;33m
NC = \033[0m # No Color

.PHONY: help
help: ## Показать справку
	@echo "$(GREEN)Доступные команды:$(NC)"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(YELLOW)  make %-20s$(NC) %s\n", $$1, $$2}'

######## Tests

.PHONY: test-cover
test-cover: ## Запустить проверку покрытия тестами
	@echo "$(GREEN)=== Running tests ===$(NC)"
	@${GO_TEST} -cover
	@echo "$(GREEN)✅ Tests completed$(NC)"

.PHONY: test
test: ## Запустить локальные тесты
	@echo "$(GREEN)=== Running tests ===$(NC)"
	@${GO_GENERATE}
	@${GO_TEST}
	@echo "$(GREEN)✅ Tests completed$(NC)"

.PHONY: test-c
test-c: ## Запустить тесты без кэширования
	@echo "$(GREEN)=== Running tests (verbose) ===$(NC)"
	@${GO_GENERATE}
	@${GO_TEST} -count=1
	@echo "$(GREEN)✅ Tests completed$(NC)"

.PHONY: test-v
test-v: ## Запустить тесты с подробным выводом
	@echo "$(GREEN)=== Running tests (verbose) ===$(NC)"
	@${GO_GENERATE}
	@${GO_TEST} -v
	@echo "$(GREEN)✅ Tests completed$(NC)"

.PHONY: test-vс
test-vc: ## Запустить тесты с подробным выводом без кэширования
	@echo "$(GREEN)=== Running tests (verbose) ===$(NC)"
	@${GO_GENERATE}
	@${GO_TEST} -v -count=1
	@echo "$(GREEN)✅ Tests completed$(NC)"

.PHONY: test-vr
test-vr: ## Запустить тесты с подробным выводом без кэширования с проверкой race detection
	@echo "$(GREEN)=== Running tests (race detection) ===$(NC)"
	@${GO_GENERATE}
	@${GO_TEST} -v -race
	@echo "$(GREEN)✅ Tests completed$(NC)"

.PHONY: test-p
test-p: ## Запустить тесты практикума
	@echo "$(GREEN)=== Running tests (practicum) ===$(NC)"
	@go build -o ./cmd/gophermart/gophermart ./cmd/gophermart/*.go \
	&& ./gophermarttest-darwin-amd64 \
                   -test.v -test.run=^TestGophermart$ \
                   -gophermart-binary-path=cmd/gophermart/gophermart \
                   -gophermart-host=localhost \
                   -gophermart-port=8080 \
                   -gophermart-database-uri="${DB_DSN}" \
                   -accrual-binary-path=cmd/accrual/accrual_darwin_amd64 \
                   -accrual-host=localhost \
                   -accrual-port=8081 \
                   -accrual-database-uri="${DB_DSN}"
	@echo "$(GREEN)✅ Tests completed$(NC)"


######## Docker

.PHONY: down
down: ## Остановить кластер
	@echo "$(YELLOW)=== Stopping cluster ===$(NC)"
	@$(DOCKER_COMPOSE) down
	@echo "$(GREEN)✅ Cluster stopped$(NC)"

.PHONY: clean
clean: ## Остановить кластер и удалить данные (volumes)
	@echo "$(YELLOW)=== Cleaning cluster data ===$(NC)"
	@$(DOCKER_COMPOSE) down -v
	@echo "$(GREEN)✅ Cluster data cleaned$(NC)"

.PHONY: up
up: ## Запустить кластер
	@echo "$(GREEN)=== Starting cluster ===$(NC)"
	@$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)✅ Cluster started$(NC)"
	@make status

.PHONY: status
status: ## Проверить статус контейнеров
	@echo "$(GREEN)=== Cluster status ===$(NC)"
	@$(DOCKER_COMPOSE) ps

.PHONY: build
build:  ## Собрать кластер
	@make down
	@make clean
	@make up
	@make status
	@sleep 5
	@make migrate-up
	@make migrate-upt
	@make unset-env

.PHONY: rebuild
rebuild:  ## пересобрать кластер
	@make down
	@make clean
	@$(DOCKER_COMPOSE) up -d --no-deps --build
	@make status
	@sleep 5
	@make migrate-up
	@make migrate-upt
	@make unset-env


######## migrate db schema
MIGRATIONS_DIR := ./migrations

.PHONY: migrate-c
migrate-c: ## Создать миграцию (make migrate-c name=[название миграции])
	@migrate create -dir=${MIGRATIONS_DIR} -ext=sql ${name}

.PHONY: migrate-up
migrate-up: ## Применить все миграции
	@migrate -path=$(MIGRATIONS_DIR) -database=$(DB_DSN) up

.PHONY: migrate-down
migrate-down: ## Откатить одну миграцию
	@migrate -path=$(MIGRATIONS_DIR) -database=$(DB_DSN) down 1


.PHONY: migrate-upt
migrate-upt: ## Применить все миграции на тестовой БД
	@migrate -path=$(MIGRATIONS_DIR) -database=$(DB_DSN_TEST) up

.PHONY: migrate-downt
migrate-downt: ## Откатить одну миграцию на тестовой БД
	@migrate -path=$(MIGRATIONS_DIR) -database=$(DB_DSN_TEST) down 1


.PHONY: show-dsn
show-dsn: ## показать DSN
	@echo DATABASE_DSN: ${DB_DSN}
	@echo DATABASE_DSN_TEST: ${DB_DSN_TEST}

.DEFAULT_GOAL := help