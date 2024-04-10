# go-go-api

Uma aplicação Go para gerenciar tenants, implementada com GORM para interação com banco de dados PostgreSQL e estruturada com práticas recomendadas para facilitar a manutenção e a escalabilidade.

## Começando

Estas instruções fornecerão uma cópia do projeto em execução na sua máquina local para fins de desenvolvimento e teste.

### Pré-requisitos

O que você precisa para instalar o software e como instalá-los:

- Go (versão 1.16 ou superior)
- PostgreSQL

### Instalação

Um passo a passo que informa o que você deve executar para ter um ambiente de desenvolvimento rodando:

1. Clone o repositório:

```bash
git clone https://github.com/jeancarlosdanese/go-go-api
```

2. Navegue até o diretório do projeto:

```bash
cd go-go-api
```

3. Instale as dependências do Go (assegure-se de que está no diretório do projeto):

```bash
go mod tidy
```

4. Crie um arquivo .env basgoo no exemplo .env.example fornecido e ajuste as configurações do banco de dados conforme necessário.

5. Execute as migrações do banco de dados:

```bash
migrate -path ./migrations -database "postgresql://postgres:postgres@localhost:5432/go_go_api?sslmode=disable" up
```

6. Execute a aplicação:

```bash
go run ./cmd/go_api
```

### Rodando os testes

```bash
go test ./...
```
