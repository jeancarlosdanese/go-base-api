# Makefile para Go Base API
.PHONY: help build run test clean docs swagger migrate up down deps install-swag

# Cores para output
GREEN := \033[0;32m
BLUE := \033[0;34m
YELLOW := \033[1;33m
NC := \033[0m # No Color

# Variáveis
APP_NAME := go-base-api
MAIN_PATH := ./cmd/go_api
BUILD_DIR := ./bin
SWAGGER_DIR := ./docs

help: ## Mostra esta ajuda
	@echo "$(BLUE)🚀 Go Base API - Comandos disponíveis:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

build: ## Compila a aplicação
	@echo "$(BLUE)🔨 Compilando aplicação...$(NC)"
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✅ Build concluído: $(BUILD_DIR)/$(APP_NAME)$(NC)"

run: ## Executa a aplicação
	@echo "$(BLUE)🚀 Executando aplicação...$(NC)"
	go run $(MAIN_PATH)/main.go

test: ## Executa os testes
	@echo "$(BLUE)🧪 Executando testes...$(NC)"
	go test -v ./...

clean: ## Limpa arquivos de build
	@echo "$(BLUE)🧹 Limpando arquivos...$(NC)"
	rm -rf $(BUILD_DIR)
	@echo "$(GREEN)✅ Arquivos limpos$(NC)"

docs: swagger ## Alias para swagger

swagger: ## Gera documentação Swagger
	@echo "$(BLUE)📝 Gerando documentação Swagger...$(NC)"
	./scripts/update_docs.sh

install-swag: ## Instala a ferramenta swag
	@echo "$(BLUE)📦 Instalando swag...$(NC)"
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "$(GREEN)✅ swag instalado$(NC)"

deps: ## Baixa as dependências
	@echo "$(BLUE)📦 Baixando dependências...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)✅ Dependências atualizadas$(NC)"

migrate-up: ## Executa migrações do banco (up)
	@echo "$(BLUE)🗄️ Executando migrações...$(NC)"
	@echo "$(YELLOW)⚠️  Certifique-se de que o banco está rodando$(NC)"
	go run $(MAIN_PATH)/main.go migrate up

migrate-down: ## Reverte última migração (down)
	@echo "$(BLUE)⬇️ Revertendo migração...$(NC)"
	go run $(MAIN_PATH)/main.go migrate down 1

lint: ## Executa linter
	@echo "$(BLUE)🔍 Executando linter...$(NC)"
	golangci-lint run

fmt: ## Formata o código
	@echo "$(BLUE)💅 Formatando código...$(NC)"
	go fmt ./...
	gofmt -s -w .

dev: ## Executa em modo desenvolvimento com hot reload
	@echo "$(BLUE)🔄 Executando em modo desenvolvimento...$(NC)"
	@echo "$(YELLOW)💡 Use 'air' para hot reload automático$(NC)"
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "$(YELLOW)💡 Para hot reload, instale: go install github.com/cosmtrek/air@latest$(NC)"; \
		go run $(MAIN_PATH)/main.go; \
	fi

docker-build: ## Build da imagem Docker
	@echo "$(BLUE)🐳 Construindo imagem Docker...$(NC)"
	docker build -t $(APP_NAME) .

docker-run: ## Executa container Docker
	@echo "$(BLUE)🐳 Executando container Docker...$(NC)"
	docker run -p 5001:5001 $(APP_NAME)

# Comandos de desenvolvimento rápidos
b: build ## Alias para build
r: run ## Alias para run
t: test ## Alias para test
d: docs ## Alias para docs

# Comando para gerar coleção do Postman a partir do Swagger
postman: ## Gera coleção do Postman a partir do Swagger
	@echo "$(BLUE)📮 Gerando coleção do Postman...$(NC)"
	@if command -v openapi2postmanv2 >/dev/null 2>&1; then \
		echo "$(YELLOW)📁 Usando arquivo local docs/swagger.json...$(NC)"; \
		if [ -f docs/swagger.json ]; then \
			echo "$(BLUE)🔄 Convertendo para formato Postman...$(NC)"; \
			openapi2postmanv2 -s docs/swagger.json -o docs/postman_collection.json; \
			echo "$(GREEN)✅ Coleção Postman gerada: postman_collection.json$(NC)"; \
			echo "$(YELLOW)💡 Importe no Postman: File > Import > postman_collection.json$(NC)"; \
		else \
			echo "$(YELLOW)⚠️  Arquivo docs/swagger.json não encontrado$(NC)"; \
			echo "$(BLUE)💡 Execute 'make docs' primeiro para gerar a documentação$(NC)"; \
		fi \
	else \
		echo "$(YELLOW)⚠️  openapi2postmanv2 não instalado$(NC)"; \
		echo "$(BLUE)💡 Instale com: npm install -g openapi-to-postmanv2$(NC)"; \
		echo "$(BLUE)💡 Ou importe diretamente via URL no Postman$(NC)"; \
		echo "$(BLUE)   URL: http://localhost:5001/swagger/doc.json$(NC)"; \
	fi

# Comando para preparar arquivos para importação no Insomnia
insomnia: ## Prepara arquivos para importação no Insomnia
	@echo "$(BLUE)🌙 Preparando arquivos para Insomnia...$(NC)"
	./scripts/generate_insomnia.sh
	@echo "$(GREEN)✅ Arquivo YAML para Insomnia criado: docs/insomnia_collection.yaml$(NC)"
	@echo "$(YELLOW)💡 Importe no Insomnia: File > Import > docs/insomnia_collection.yaml$(NC)"
