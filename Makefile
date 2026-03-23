.PHONY: help build run test clean migrate-up migrate-down dev
include .env
export
# Переменные
APP_NAME := crm-api
MAIN_PATH := ./cmd/api
BUILD_DIR := ./bin

# Цвета для вывода
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m # No Color

help: ## Показать эту справку
	@echo "$(GREEN)=== $(APP_NAME) - Available Commands ===$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "$(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

build: ## Скомпилировать приложение
	@echo "$(GREEN)Building $(APP_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Build completed!$(NC)"

run: ## Запустить приложение
	@echo "$(GREEN)Starting $(APP_NAME)...$(NC)"
	@go run $(MAIN_PATH)/main.go

dev: ## Запустить с air (hot reload)
	@echo "$(GREEN)Starting with hot reload...$(NC)"
	@air

test: ## Запустить тесты
	@echo "$(GREEN)Running tests...$(NC)"
	@go test -v ./...

test-coverage: ## Запустить тесты с покрытием
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

clean: ## Очистить сборку
	@echo "$(YELLOW)Cleaning...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)Cleaned!$(NC)"

migrate-create: ##Создать миграцию
	@echo "$(GREEN)Running migrations file create...$(NC)"
	@migrate create -ext sql -dir migrations -seq ${name}

migrate-up: ## Применить миграции вверх
	@echo "$(GREEN)Running migrations up...$(NC)"
	@migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

migrate-down: ## Откатить миграции вниз
	@echo "$(YELLOW)Rolling back migrations...$(NC)"
	@migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down


migrate-force:
	@echo "$(YELLOW)Forcing migration version to $(VERSION)...$(NC)"
	@migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" force $(VERSION)

install-deps: ## Установить все зависимости
	@echo "$(GREEN)Installing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Dependencies installed!$(NC)"

lint: ## Запустить линтер
	@echo "$(GREEN)Running linter...$(NC)"
	@golangci-lint run ./...

fmt: ## Форматировать код
	@echo "$(GREEN)Formatting code...$(NC)"
	@go fmt ./...

vet: ## Запустить go vet
	@echo "$(GREEN)Running go vet...$(NC)"
	@go vet ./...

# Комбинации
pre-commit: fmt vet lint ## Проверки перед коммитом
	@echo "$(GREEN)All checks passed!$(NC)"

all: clean build test ## Очистить, собрать, протестировать
	@echo "$(GREEN)All done!$(NC)"