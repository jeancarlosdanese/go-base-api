# @file: Makefile
# Makefile para Go Base API - Numerologia Cabalística
.PHONY: help build run test clean docs swagger migrate-up migrate-down deps install-swag lint fmt dev docker-build docker-run postman insomnia

# Cores para output
GREEN := \033[0;32m
BLUE := \033[0;34m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Variáveis
APP_NAME := essentia-api
MAIN_PATH := ./cmd/go_api
BUILD_DIR := ./bin
SWAGGER_DIR := ./docs
PORT := 5001

help: ## Mostra esta ajuda
	@echo "$(BLUE)🔮 Go Base API - Comandos disponíveis:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

build: ## Compila a aplicação
	@echo "$(BLUE)🔨 Compilando aplicação...$(NC)"
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✅ Build concluído: $(BUILD_DIR)/$(APP_NAME)$(NC)"

run: ## Executa a aplicação
	@echo "$(BLUE)🚀 Executando aplicação na porta $(PORT)...$(NC)"
	@echo "$(YELLOW)💡 Certifique-se de que o PostgreSQL e Redis estão rodando$(NC)"
	go run $(MAIN_PATH)/main.go

test: ## Executa os testes
	@echo "$(BLUE)🧪 Executando testes...$(NC)"
	go test -v ./...

test-coverage: ## Executa testes com cobertura
	@echo "$(BLUE)🧪 Executando testes com cobertura...$(NC)"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Relatório de cobertura gerado: coverage.html$(NC)"

clean: ## Limpa arquivos de build e cache
	@echo "$(BLUE)🧹 Limpando arquivos...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	go clean -cache -testcache
	@echo "$(GREEN)✅ Arquivos limpos$(NC)"

docs: swagger ## Alias para swagger

swagger: ## Gera documentação Swagger
	@echo "$(BLUE)📝 Gerando documentação Swagger...$(NC)"
	@if [ -f "./scripts/update_docs.sh" ]; then \
		./scripts/update_docs.sh; \
	else \
		echo "$(YELLOW)⚠️  Script update_docs.sh não encontrado$(NC)"; \
		echo "$(BLUE)💡 Executando swag init diretamente...$(NC)"; \
		swag init -g $(MAIN_PATH)/main.go --output $(SWAGGER_DIR)/ --parseDependency --parseInternal --parseDepth 2 --useStructName; \
	fi
	@echo "$(GREEN)✅ Documentação Swagger gerada$(NC)"

install-swag: ## Instala a ferramenta swag
	@echo "$(BLUE)📦 Instalando swag...$(NC)"
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "$(GREEN)✅ swag instalado$(NC)"

deps: ## Baixa e atualiza as dependências
	@echo "$(BLUE)📦 Baixando dependências...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)✅ Dependências atualizadas$(NC)"

# Função auxiliar para obter DATABASE_URL do .env ou construir a partir das variáveis individuais
get-database-url = $(shell \
	if [ -f .env ]; then \
		DB_URL=$$(grep "^DATABASE_URL=" .env 2>/dev/null | cut -d '=' -f2- | sed 's/^"//;s/"$$//'); \
		if [ -n "$$DB_URL" ]; then \
			echo "$$DB_URL"; \
		else \
			DB_USER=$$(grep "^DB_USER=" .env 2>/dev/null | cut -d '=' -f2- | sed 's/^"//;s/"$$//' | head -1 || echo "hyberica"); \
			DB_PASSWORD=$$(grep "^DB_PASSWORD=" .env 2>/dev/null | cut -d '=' -f2- | sed 's/^"//;s/"$$//' | head -1 || echo "hyberica"); \
			DB_HOST=$$(grep "^DB_HOST=" .env 2>/dev/null | cut -d '=' -f2- | sed 's/^"//;s/"$$//' | head -1 || echo "localhost"); \
			DB_PORT=$$(grep "^DB_PORT=" .env 2>/dev/null | cut -d '=' -f2- | sed 's/^"//;s/"$$//' | head -1 || echo "5432"); \
			DB_NAME=$$(grep "^DB_NAME=" .env 2>/dev/null | cut -d '=' -f2- | sed 's/^"//;s/"$$//' | head -1 || echo "go_api"); \
			echo "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable"; \
		fi; \
	else \
		echo "postgres://hyberica:hyberica@localhost:5432/go_api?sslmode=disable"; \
	fi \
)

migrate-up: ## Executa migrações do banco (up)
	@echo "$(BLUE)🗄️ Executando migrações...$(NC)"
	@echo "$(YELLOW)⚠️  Certifique-se de que o PostgreSQL está rodando$(NC)"
	@if command -v migrate >/dev/null 2>&1; then \
		DB_URL=$(call get-database-url); \
		echo "$(BLUE)📊 Database URL: postgres://$$(echo $$DB_URL | sed 's/:[^@]*@/:***@/' | sed 's/@[^/]*\//@***\//')$(NC)"; \
		migrate -path ./migrations -database "$$DB_URL" up 1; \
	else \
		echo "$(YELLOW)⚠️  migrate não instalado. Instale com: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest$(NC)"; \
	fi

migrate-down: ## Reverte última migração (down)
	@echo "$(BLUE)⬇️ Revertendo migração...$(NC)"
	@if command -v migrate >/dev/null 2>&1; then \
		DB_URL=$(call get-database-url); \
		echo "$(BLUE)📊 Database URL: postgres://$$(echo $$DB_URL | sed 's/:[^@]*@/:***@/' | sed 's/@[^/]*\//@***\//')$(NC)"; \
		migrate -path ./migrations -database "$$DB_URL" down 1; \
	else \
		echo "$(YELLOW)⚠️  migrate não instalado. Instale com: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest$(NC)"; \
	fi

lint: ## Executa linter (requer golangci-lint)
	@echo "$(BLUE)🔍 Executando linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint não instalado$(NC)"; \
		echo "$(BLUE)💡 Instale com: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

fmt: ## Formata o código
	@echo "$(BLUE)💅 Formatando código...$(NC)"
	go fmt ./...
	@echo "$(GREEN)✅ Código formatado$(NC)"

vet: ## Executa go vet para análise estática
	@echo "$(BLUE)🔍 Executando go vet...$(NC)"
	go vet ./...
	@echo "$(GREEN)✅ Análise estática concluída$(NC)"

dev: ## Executa em modo desenvolvimento com hot reload
	@echo "$(BLUE)🔄 Executando em modo desenvolvimento...$(NC)"
	@echo "$(YELLOW)💡 Use 'air' para hot reload automático$(NC)"
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "$(YELLOW)💡 Para hot reload, instale: go install github.com/cosmtrek/air@latest$(NC)"; \
		echo "$(BLUE)🚀 Executando sem hot reload...$(NC)"; \
		go run $(MAIN_PATH)/main.go; \
	fi

docker-build: ## Build da imagem Docker
	@echo "$(BLUE)🐳 Construindo imagem Docker...$(NC)"
	@if [ -f "Dockerfile" ]; then \
		docker build -t $(APP_NAME):latest .; \
		echo "$(GREEN)✅ Imagem Docker construída: $(APP_NAME):latest$(NC)"; \
	else \
		echo "$(RED)❌ Dockerfile não encontrado$(NC)"; \
	fi

docker-run: ## Executa container Docker
	@echo "$(BLUE)🐳 Executando container Docker...$(NC)"
	docker run --rm -p $(PORT):$(PORT) --env-file .env $(APP_NAME):latest

# Comandos de desenvolvimento rápidos (aliases)
b: build ## Alias para build
r: run ## Alias para run
t: test ## Alias para test
d: docs ## Alias para docs
c: clean ## Alias para clean

# Comando para gerar coleção do Postman a partir do Swagger
postman: ## Gera coleção do Postman a partir do Swagger
	@echo "$(BLUE)📮 Gerando coleção do Postman...$(NC)"
	@if command -v openapi2postmanv2 >/dev/null 2>&1; then \
		echo "$(YELLOW)📁 Usando arquivo local docs/swagger.json...$(NC)"; \
		if [ -f docs/swagger.json ]; then \
			echo "$(BLUE)🔄 Convertendo para formato Postman...$(NC)"; \
			openapi2postmanv2 -s docs/swagger.json -o docs/postman_collection.json; \
			echo "$(GREEN)✅ Coleção Postman gerada: docs/postman_collection.json$(NC)"; \
			echo "$(YELLOW)💡 Importe no Postman: File > Import > docs/postman_collection.json$(NC)"; \
		else \
			echo "$(YELLOW)⚠️  Arquivo docs/swagger.json não encontrado$(NC)"; \
			echo "$(BLUE)💡 Execute 'make docs' primeiro para gerar a documentação$(NC)"; \
		fi \
	else \
		echo "$(YELLOW)⚠️  openapi2postmanv2 não instalado$(NC)"; \
		echo "$(BLUE)💡 Instale com: npm install -g openapi-to-postmanv2$(NC)"; \
		echo "$(BLUE)💡 Ou importe diretamente via URL no Postman:$(NC)"; \
		echo "$(BLUE)   URL: http://localhost:$(PORT)/swagger/doc.json$(NC)"; \
	fi

# Comando para preparar arquivos para importação no Insomnia
insomnia: ## Prepara arquivos para importação no Insomnia
	@echo "$(BLUE)🌙 Preparando arquivos para Insomnia...$(NC)"
	@if [ -f "./scripts/generate_insomnia.sh" ]; then \
		./scripts/generate_insomnia.sh; \
		echo "$(GREEN)✅ Arquivo YAML para Insomnia criado: docs/insomnia_collection.yaml$(NC)"; \
		echo "$(YELLOW)💡 Importe no Insomnia: File > Import > docs/insomnia_collection.yaml$(NC)"; \
	else \
		echo "$(YELLOW)⚠️  Script generate_insomnia.sh não encontrado$(NC)"; \
	fi

# Comandos úteis
check: fmt vet test ## Formata, analisa e testa o código
	@echo "$(GREEN)✅ Verificação completa concluída$(NC)"

version: ## Mostra versão da aplicação
	@echo "$(BLUE)📦 Versão:$(NC)"
	@if [ -f "VERSION" ]; then \
		cat VERSION; \
	else \
		echo "Desenvolvimento"; \
	fi
