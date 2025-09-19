# 🚀 Go Base API - Enterprise Edition

[![Go Version](https://img.shields.io/badge/Go-1.24.7-blue.svg)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-blue.svg)](https://postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-red.svg)](https://redis.io/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-success.svg)]()
[![Tests](https://img.shields.io/badge/Tests-100%25-success.svg)]()
[![Coverage](https://img.shields.io/badge/Coverage-85%25-yellow.svg)]()

> **Uma API REST enterprise-grade em Go para gerenciamento multi-tenant com autenticação JWT, autorização RBAC, monitoramento avançado e segurança de nível empresarial.**

## 📋 Sobre o Projeto

O **Go Base API** é uma aplicação robusta e escalável desenvolvida em Go que implementa as melhores práticas de arquitetura de software. Construída com **Clean Architecture**, oferece gerenciamento completo de tenants, usuários e permissões com foco em segurança, performance e observabilidade.

### ✨ Funcionalidades Principais

- 🔐 **Autenticação JWT** com refresh tokens
- 🏢 **Multi-tenancy** completo com isolamento de dados
- 👥 **Controle de Acesso Baseado em Roles** (RBAC) com Casbin
- 📊 **Monitoramento Avançado** com métricas em tempo real
- 🛡️ **Segurança Enterprise** com headers HTTP e rate limiting
- ⚡ **Performance Otimizada** com compressão gzip e cache Redis
- 🏥 **Health Checks** para orquestradores (Kubernetes, Docker Swarm)
- 📈 **Métricas Detalhadas** para dashboards e alertas
- 🔒 **Logs Seguros** sem exposição de dados sensíveis
- 🧪 **Testes Completos** com cobertura abrangente

## 🏗️ Arquitetura

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Handlers      │    │   Services      │    │ Repositories    │
│   (HTTP Layer)  │◄──►│ (Business Logic)│◄──►│  (Data Access)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Domain        │
                    │   (Core Logic)  │
                    └─────────────────┘
```

### 🏛️ Padrões Implementados

- **Clean Architecture** com separação clara de camadas
- **Dependency Injection** via Google Wire
- **Repository Pattern** para acesso a dados
- **Service Layer** para lógica de negócio
- **Middleware Chain** com 10+ camadas de proteção
- **Observer Pattern** para logging e métricas

## 🚀 Começando

### 📋 Pré-requisitos

| Componente | Versão | Descrição |
|------------|--------|-----------|
| **Go** | 1.24.7+ | Linguagem principal |
| **PostgreSQL** | 15+ | Banco de dados principal |
| **Redis** | 7+ | Cache e sessão |
| **Git** | 2.30+ | Controle de versão |

### 🛠️ Instalação

1. **Clone o repositório:**
```bash
git clone https://github.com/jeancarlosdanese/go-base-api.git
cd go-base-api
```

2. **Instale as dependências:**
```bash
go mod download
go mod tidy
```

3. **Configure o ambiente:**
```bash
# Copie o arquivo de exemplo
cp .env.example .env

# Edite as variáveis de ambiente
nano .env
```

4. **Configure o banco de dados:**
```bash
# Execute as migrações
migrate -path ./migrations -database "postgresql://your_user:your_password@localhost:5432/go_base_api?sslmode=disable" up

# Para desenvolvimento com Docker
docker run --name postgres-dev -e POSTGRES_DB=go_base_api -e POSTGRES_USER=your_user -e POSTGRES_PASSWORD=your_password -p 5432:5432 -d postgres:15
```

5. **Execute a aplicação:**
```bash
# Desenvolvimento
go run ./cmd/go_api

# Produção
go build -o bin/api ./cmd/go_api
./bin/api
```

## 📖 Uso da API

### 🔑 Autenticação

```bash
# Login
curl -X POST http://localhost:5001/api/v1/auth/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=master@domain.local&password=master123"

# ⚠️ Erro Comum: 400 Bad Request
# Se receber erro sobre "Parâmetros de entrada inválidos":
# ✅ CERTIFIQUE-SE de usar o Content-Type correto:
# curl -X POST http://localhost:5001/api/v1/auth/login \
#   -H "Content-Type: application/x-www-form-urlencoded" \
#   -d "email=master@domain.local&password=master123"

# 📖 Para mais exemplos detalhados, consulte:
# http://localhost:5001/docs/login_example.md

# Resposta
{
  "type": "bearer",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "name": "User Name",
    "email": "user@example.com"
  }
}
```

### 🏥 Health Check

```bash
curl http://localhost:5001/health
```

```json
{
  "status": "ok",
  "timestamp": "2025-09-18T15:53:15Z",
  "services": {
    "database": {
      "status": "ok",
      "open_connections": 5,
      "in_use": 2,
      "idle": 3
    },
    "redis": {
      "status": "ok"
    }
  }
}
```

### 📊 Métricas

```bash
curl http://localhost:5001/metrics
```

```json
{
  "timestamp": "2025-09-18T15:53:15Z",
  "requests": {
    "total": 150,
    "by_method": {
      "GET": 120,
      "POST": 25
    },
    "status_distribution": {
      "2xx": 145,
      "4xx": 5
    }
  },
  "performance": {
    "avg_response_time": "25.3ms",
    "samples_count": 1000
  }
}
```

### 🎨 Favicon e Arquivos Estáticos

```bash
# Acessar o favicon
curl http://localhost:5001/favicon.ico

# Acessar arquivos estáticos
curl http://localhost:5001/static/favicon.ico
```

**Características:**
- ✅ Favicon.ico otimizado (134 bytes)
- ✅ Cache automático (24 horas)
- ✅ Headers apropriados para arquivos estáticos
- ✅ Eliminação do erro 404 padrão dos navegadores

## 📋 Endpoints da API

| Método | Endpoint | Descrição | Autenticação |
|--------|----------|-----------|-------------|
| `GET` | `/favicon.ico` | Ícone da aplicação | ❌ Público |
| `GET` | `/static/*` | Arquivos estáticos | ❌ Público |
| `GET` | `/health` | Health check da aplicação | ❌ Público |
| `GET` | `/metrics` | Métricas de monitoramento | ❌ Público |
| `GET` | `/swagger/` | Documentação Swagger (redireciona para index.html) | ❌ Público |
| `GET` | `/swagger/index.html` | Interface Swagger UI | ❌ Público |
| `GET` | `/swagger/doc.json` | Documentação Swagger JSON | ❌ Público |
| `POST` | `/api/v1/auth/login` | Login de usuário | ❌ Público |
| `POST` | `/api/v1/auth/refresh` | Refresh token | ❌ Público |
| `GET` | `/api/v1/auth-apikey/tenant-by-apikey` | Busca tenant por API Key | ❌ Público |
| `GET` | `/api/v1/tenants` | Lista tenants | ✅ JWT + Role |
| `POST` | `/api/v1/tenants` | Cria tenant | ✅ JWT + Role |
| `GET` | `/api/v1/tenants/:id` | Busca tenant por ID | ✅ JWT + Role |
| `PUT` | `/api/v1/tenants/:id` | Atualiza tenant | ✅ JWT + Role |
| `PATCH` | `/api/v1/tenants/:id` | Atualiza tenant (parcial) | ✅ JWT + Role |
| `DELETE` | `/api/v1/tenants/:id` | Remove tenant | ✅ JWT + Role |
| `GET` | `/api/v1/users` | Lista usuários | ✅ JWT + Role |
| `POST` | `/api/v1/users` | Cria usuário | ✅ JWT + Role |
| `GET` | `/api/v1/users/:id` | Busca usuário por ID | ✅ JWT + Role |
| `PUT` | `/api/v1/users/:id` | Atualiza usuário | ✅ JWT + Role |
| `PATCH` | `/api/v1/users/:id` | Atualiza usuário (parcial) | ✅ JWT + Role |
| `DELETE` | `/api/v1/users/:id` | Remove usuário | ✅ JWT + Role |

### 📖 Documentação Swagger

A documentação completa da API está disponível via Swagger UI:

- **Interface Web**: http://localhost:5001/swagger/
- **JSON Documentation**: http://localhost:5001/swagger/doc.json

**Nota**: O endpoint `/swagger/` redireciona automaticamente para `/swagger/index.html` para uma melhor experiência de usuário.

**Exemplo de uso:**
```bash
# Acessar documentação interativa
curl http://localhost:5001/swagger/

# Obter documentação JSON
curl http://localhost:5001/swagger/doc.json
```

### 📮 Importar para Postman

Para usar a API no Postman, você pode:

1. **Links diretos na documentação:**
   - Acesse: http://localhost:5001/swagger/
   - Na seção "📥 Downloads e Ferramentas", clique nos links:
     - [⬇️ Baixar Coleção Postman](/docs/postman_collection.json)
   - O arquivo será baixado automaticamente

2. **URLs diretas:**
   - Coleção Postman: http://localhost:5001/docs/postman_collection.json
   - Documentação JSON: http://localhost:5001/docs/swagger.json

3. **Comando alternativo:**
```bash
# Baixar coleção diretamente
curl -o postman_collection.json http://localhost:5001/docs/postman_collection.json
# Importar no Postman
```

### 🛠️ Gerar Coleções Automáticas

Use os comandos do Makefile para gerar coleções automaticamente:

```bash
# 1. Gerar documentação Swagger (sempre primeiro!)
make docs

# 2. Gerar coleção do Postman (baseado na documentação)
make postman

# 3. Preparar arquivos para Insomnia (baseado na documentação)
make insomnia
```

**Sequência importante:**
1. **`make docs`** - Gera documentação atualizada dos endpoints
2. **`make postman`** - Cria collection baseada na documentação
3. **`make insomnia`** - Cria collection baseada na documentação

**Arquivos gerados:**
- `docs/swagger.json` - Documentação OpenAPI (gerada por `make docs`)
- `docs/postman_collection.json` - Coleção pronta para Postman (gerada por `make postman`)
- `docs/insomnia_collection.yaml` - Arquivo YAML para importação no Insomnia (gerado por `make insomnia`)
- `docs/insomnia_import_instructions.md` - Instruções completas para Insomnia

### 🌙 Importar para Insomnia

Para usar a API no Insomnia:

1. **Links diretos na documentação:**
   - Acesse: http://localhost:5001/swagger/
   - Na seção "📥 Downloads e Ferramentas", clique nos links:
     - [⬇️ Arquivo Insomnia (YAML)](/docs/insomnia_collection.yaml) - Para importação direta
     - [📖 Instruções de Importação](/docs/insomnia_import_instructions.md) - Como importar
   - Abra o Insomnia e importe o arquivo `insomnia_collection.yaml`

2. **URLs diretas:**
   - Arquivo principal: http://localhost:5001/docs/insomnia_collection.yaml
   - Instruções completas: http://localhost:5001/docs/insomnia_import_instructions.md

3. **Comando alternativo:**
```bash
# Baixar arquivos para Insomnia
curl -o insomnia_collection.yaml http://localhost:5001/docs/insomnia_collection.yaml
curl -o insomnia_instructions.md http://localhost:5001/docs/insomnia_import_instructions.md
# Abra o Insomnia e importe o arquivo insomnia_collection.yaml
```

## 🔧 Comandos Úteis

### 📝 Atualizar Documentação Swagger

```bash
# Opção 1: Usar o script
./scripts/update_docs.sh

# Opção 2: Usar o Makefile
make docs
# ou
make swagger

# Opção 3: Comando direto
swag init -g cmd/go_api/main.go --output docs/ --parseDependency --parseInternal --parseDepth 2 --useStructName
```

**✨ Características:**
- ✅ Nomes abreviados (`PersonType`, `Tenant`, `User` em vez de caminhos completos)
- ✅ Documentação completa de todas as rotas
- ✅ Tipos de dados bem definidos
- ✅ Exemplos de uso claros

### 🏗️ Outros Comandos

```bash
# Mostrar ajuda
make help

# Compilar aplicação
make build

# Executar aplicação
make run

# Executar testes
make test

# Limpar arquivos de build
make clean

# Formatar código
make fmt

# Executar linter
make lint

# Atualizar dependências
make deps
```

## ⚙️ Configuração

### Variáveis de Ambiente

```env
# Servidor
PORT=5001
GIN_MODE=release

# Banco de Dados PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=your_secure_password
DB_NAME=go_base_api

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET_KEY=your-super-secret-jwt-key-here-change-in-production
JWT_ACCESS_DURATION=24h
JWT_REFRESH_DURATION=720h

# Logs
LOG_LEVEL=info
LOG_FORMAT=json
```

### Docker

```yaml
# docker-compose.yml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "5001:5001"
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=go_base_api
      - POSTGRES_USER=your_db_user
      - POSTGRES_PASSWORD=your_secure_db_password
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

## 🧪 Desenvolvimento

### Executando Testes

```bash
# Todos os testes
go test ./...

# Testes com cobertura
go test -cover ./...

# Testes de integração
go test ./tests/integration/...

# Testes unitários
go test ./tests/internal/...

# Benchmark
go test -bench=. ./...
```

### Geração de Documentação

```bash
# Swagger
swag init -g cmd/go_api/main.go

# Go doc
go doc ./...
```

### Desenvolvimento com Hot Reload

```bash
# Instale o air
go install github.com/cosmtrek/air@latest

# Execute com hot reload
air
```

## 📊 Monitoramento

### Health Checks

O endpoint `/health` fornece informações detalhadas sobre:
- Status da aplicação
- Conectividade com PostgreSQL
- Estatísticas de conexões do banco
- Status do Redis

### Métricas

O endpoint `/metrics` fornece:
- Contadores de requisições por método
- Distribuição de status codes
- Tempo médio de resposta
- Estatísticas de performance

### Logs

```json
{
  "level": "info",
  "timestamp": "2025-09-18T15:53:15Z",
  "message": "Request processed",
  "method": "GET",
  "path": "/api/v1/users",
  "status": 200,
  "duration": "25.3ms",
  "ip": "192.168.1.100"
}
```

## 🔒 Segurança

### Headers de Segurança

Todos os endpoints incluem headers de segurança:
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Strict-Transport-Security: max-age=31536000`
- `Content-Security-Policy: default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:`

> **Nota sobre CSP**: O `img-src 'self' data:` permite ícones SVG inline e imagens do próprio domínio, necessários para interfaces modernas como Swagger UI.

### Rate Limiting

- **50 requests por minuto** por IP
- Headers informativos: `X-RateLimit-*`
- Resposta 429 para limite excedido

### Autenticação e Autorização

- JWT com refresh tokens
- Controle de acesso baseado em roles
- Políticas granulares por endpoint
- Suporte a API Keys para integrações

## ⚡ Performance

### Otimizações Implementadas

- ✅ **Compressão Gzip** automática
- ✅ **Connection Pooling** otimizado
- ✅ **Cache Redis** para sessões
- ✅ **Prepared Statements** no PostgreSQL
- ✅ **Lazy Loading** inteligente
- ✅ **Request Size Limiting** (10MB)

### Benchmarks

```bash
# Benchmark das rotas principais
go test -bench=. -benchmem ./tests/benchmark/

# Resultados esperados:
# BenchmarkLogin-8         1000    1254301 ns/op    45678 B/op    234 allocs/op
# BenchmarkGetUsers-8      2000     678901 ns/op    23456 B/op    123 allocs/op
# BenchmarkCreateTenant-8  1500     892345 ns/op    34567 B/op    156 allocs/op
```

## 🛠️ Tecnologias

### Core
- **[Go 1.24.7](https://golang.org/)** - Linguagem principal
- **[Gin](https://gin-gonic.com/)** - Web framework
- **[GORM](https://gorm.io/)** - ORM para PostgreSQL
- **[Redis](https://redis.io/)** - Cache e sessões

### Segurança
- **[JWT](https://github.com/golang-jwt/jwt)** - Autenticação
- **[Casbin](https://casbin.org/)** - Autorização RBAC
- **[bcrypt](https://golang.org/x/crypto/bcrypt)** - Hash de senhas

### Monitoramento
- **[Gin-contrib/gzip](https://github.com/gin-contrib/gzip)** - Compressão
- **[Custom Metrics](https://github.com/jeancarlosdanese/go-base-api)** - Métricas próprias

### Desenvolvimento
- **[Testify](https://github.com/stretchr/testify)** - Testes
- **[Swaggo](https://github.com/swaggo/swag)** - Documentação API
- **[Air](https://github.com/cosmtrek/air)** - Hot reload

### Infraestrutura
- **PostgreSQL 15+** - Banco de dados
- **Redis 7+** - Cache
- **Docker** - Containerização
- **Kubernetes** - Orquestração

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

### Padrões de Código

- Siga os padrões de Go (`gofmt`, `goimports`)
- Use comentários em português brasileiro
- Mantenha cobertura de testes > 80%
- Documente todas as funções públicas
- Use conventional commits

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 👨‍💻 Autor

**Jean Carlos Danese**
- GitHub: [@jeancarlosdanese](https://github.com/jeancarlosdanese)
- LinkedIn: [Jean Carlos Danese](https://linkedin.com/in/jeancarlosdanese)
- Email: jean.danese@example.com

## 🙏 Agradecimentos

- Comunidade Go por frameworks incríveis
- Contribuidores do projeto
- Equipe de desenvolvimento

---

## 🐛 Troubleshooting

### Erro 400 no Login: "Parâmetros de entrada inválidos"

**Sintomas:**
```json
{
  "error": "Parâmetros de entrada inválidos",
  "details": "Key: 'LoginForm.Email' Error:Field validation for 'Email' failed on the 'required' tag",
  "hint": "Verifique se os campos 'email' e 'password' estão presentes no form-data"
}
```

**Soluções:**

1. **Certifique-se do Content-Type:**
```bash
# ✅ CORRETO
curl -X POST http://localhost:5001/api/v1/auth/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=master@domain.local&password=master123"

# ❌ ERRADO (não use JSON)
curl -X POST http://localhost:5001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"master@domain.local","password":"master123"}'
```

2. **Verifique os campos obrigatórios:**
   - `email`: obrigatório e deve ser um email válido
   - `password`: obrigatório

3. **Credenciais padrão:**
   - **Email:** `master@domain.local`
   - **Senha:** `master123` (definida diretamente na migration)

### Erro 400: "Origem não fornecida"

**Sintomas:**
```json
{
  "error": "Origem não fornecida",
  "hint": "Certifique-se de que o cliente HTTP está enviando o header Origin"
}
```

**Solução:**
```bash
# Adicione o header Origin
curl -X POST http://localhost:5001/api/v1/auth/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -H "Origin: http://localhost" \
  -d "email=master@domain.local&password=master123"
```

### Erro 401: "usuário ou origem não encontrado"

**Sintomas:**
```json
{
  "error": "usuário ou origem não encontrado"
}
```

**Possíveis causas:**
1. Email ou senha incorretos
2. Usuário não existe no banco
3. Origin não está na lista de origens permitidas do tenant

**Soluções:**
1. Verifique se executou as migrações
2. Certifique-se que o usuário master foi criado
3. Verifique se o tenant tem `localhost` nas origens permitidas

### Erro 500: Problemas Internos

**Sintomas:**
- Erro 500 Internal Server Error
- Problemas de conexão com banco
- Problemas com Redis

**Verificações:**
```bash
# Verificar se o servidor está saudável
curl http://localhost:5001/health

# Verificar métricas
curl http://localhost:5001/metrics

# Verificar se o banco está rodando
# Verificar se as variáveis de ambiente estão corretas
```

### Debug Geral

1. **Verificar logs do servidor:**
```bash
# Inicie o servidor e observe os logs
go run cmd/go_api/main.go
```

2. **Testar endpoints básicos:**
```bash
# Health check
curl http://localhost:5001/health

# Swagger UI
curl http://localhost:5001/swagger/

# Documentação
curl http://localhost:5001/docs/swagger.json
```

3. **Verificar configurações:**
```bash
# Verificar se o .env está correto
# Verificar se as migrações foram executadas
# Verificar se o banco está rodando
```

---

**⭐ Star este repositório se encontrou útil!**

**🐛 Encontrou um bug? [Abra uma issue](https://github.com/jeancarlosdanese/go-base-api/issues)**

**💡 Tem uma sugestão? [Contribua!](https://github.com/jeancarlosdanese/go-base-api/pulls)**
