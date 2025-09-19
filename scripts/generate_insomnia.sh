#!/bin/bash

# Script para gerar arquivos compatíveis com Insomnia
# Cria arquivo JSON para importação direta + instruções
# Uso: ./scripts/generate_insomnia.sh

echo "🚀 Preparando arquivos para importação no Insomnia..."

# Verificar se a documentação Swagger existe
if [ ! -f "docs/swagger.json" ]; then
    echo "❌ Arquivo docs/swagger.json não encontrado!"
    echo "💡 Execute primeiro: make docs"
    exit 1
fi

# Ler informações da documentação Swagger
echo "📄 Analisando documentação Swagger..."

# Usar jq se disponível, senão usar grep como fallback
if command -v jq &> /dev/null; then
    SWAGGER_TITLE=$(jq -r '.info.title // "Go Base API"' docs/swagger.json)
    SWAGGER_VERSION=$(jq -r '.info.version // "1.0.0"' docs/swagger.json)
    SWAGGER_DESCRIPTION=$(jq -r '.info.description // "API desenvolvida em Go"' docs/swagger.json)
else
    SWAGGER_TITLE=$(grep -o '"title":"[^"]*"' docs/swagger.json | head -1 | sed 's/.*"title":"\([^"]*\)".*/\1/')
    SWAGGER_VERSION=$(grep -o '"version":"[^"]*"' docs/swagger.json | head -1 | sed 's/.*"version":"\([^"]*\)".*/\1/')
    SWAGGER_DESCRIPTION="API desenvolvida em Go"
fi

echo "📋 API: $SWAGGER_TITLE v$SWAGGER_VERSION"

# Criar arquivo YAML compatível com Collection do Insomnia
echo "📄 Criando collection do Insomnia..."

# Criar collection estruturada manualmente (mais confiável) no formato YAML
cat > docs/insomnia_collection.yaml << EOF
type: collection.insomnia.rest/5.0
name: "$SWAGGER_TITLE"
meta:
  id: wrk_8354f20c4090420684bd7ff411ab4a57
  created: 1758235337828
  modified: $(date +%s)000
  description: "$SWAGGER_DESCRIPTION - Multi-tenant com autenticação JWT (v$SWAGGER_VERSION)"
collection:
  - url: "{{ base_url }}/health"
    name: 🏥 Health Check
    meta:
      id: req_5081add2a8c2470abc5ed5465e6d67ac
      created: 1758235337836
      modified: 1758235337836
      isPrivate: false
      description: ""
      sortKey: -1758235337836
    method: GET
    settings:
      renderRequestBody: true
      encodeUrl: true
      followRedirects: global
      cookies:
        send: true
        store: true
      rebuildPath: true
  - url: "{{ base_url }}/swagger/"
    name: 📚 Documentação Swagger
    meta:
      id: req_7c9b47d66e294dd89d9db444297ed2a9
      created: 1758235337836
      modified: 1758235337836
      isPrivate: false
      description: ""
      sortKey: -1758235337836
    method: GET
    settings:
      renderRequestBody: true
      encodeUrl: true
      followRedirects: global
      cookies:
        send: true
        store: true
      rebuildPath: true
  - url: "{{ base_url }}/metrics"
    name: 📊 Métricas
    meta:
      id: req_b2b54f3ee5da4627b2ace783ae718456
      created: 1758235337836
      modified: 1758235337836
      isPrivate: false
      description: ""
      sortKey: -1758235337836
    method: GET
    settings:
      renderRequestBody: true
      encodeUrl: true
      followRedirects: global
      cookies:
        send: true
        store: true
      rebuildPath: true
  - name: 🔐 Autenticação
    meta:
      id: fld_1a13be8cfd484b6abac68adff7eabf62
      created: 1758235337837
      modified: 1758235337837
      sortKey: -1758235337837
      description: ""
    children:
      - url: "{{ base_url }}/api/v1/auth/login"
        name: 🔐 Login - Obter Token JWT
        meta:
          id: req_fecbadceb63348f0b836795ac259f662
          created: 1758235337831
          modified: 1758235337831
          isPrivate: false
          description: Faça login com email e senha para obter token de acesso JWT
          sortKey: -1758235337831
        method: POST
        body:
          mimeType: application/x-www-form-urlencoded
          params:
            - name: email
              value: "{{ master_email }}"
            - name: password
              value: "{{ master_password }}"
        headers:
          - name: Content-Type
            value: application/x-www-form-urlencoded
          - name: Origin
            value: "{{ origin }}"
          - name: User-Agent
            value: "{{ user_agent }}"
          - name: Accept
            value: application/json
        settings:
          renderRequestBody: true
          encodeUrl: true
          followRedirects: global
          cookies:
            send: true
            store: true
          rebuildPath: true
      - url: "{{ base_url }}/api/v1/auth/refresh"
        name: Refresh Token
        meta:
          id: req_04eac6dd00de4a3fb64c83b170413a06
          created: 1758235337833
          modified: 1758235337833
          isPrivate: false
          description: ""
          sortKey: -1758235337833
        method: POST
        body:
          mimeType: application/x-www-form-urlencoded
          params:
            - name: refreshToken
              value: "{{ refresh_token }}"
        settings:
          renderRequestBody: true
          encodeUrl: true
          followRedirects: global
          cookies:
            send: true
            store: true
          rebuildPath: true
      - url: "{{ base_url }}/api/v1/auth-apikey/tenant-by-apikey"
        name: Buscar Tenant por API Key
        meta:
          id: req_23b5b9e200f044d19a8302f56ff0d844
          created: 1758235337834
          modified: 1758235337834
          isPrivate: false
          description: ""
          sortKey: -1758235337834
        method: GET
        headers:
          - name: X-API-Key
            value: "{{ api_key }}"
        settings:
          renderRequestBody: true
          encodeUrl: true
          followRedirects: global
          cookies:
            send: true
            store: true
          rebuildPath: true
  - name: 👥 Usuários
    meta:
      id: fld_e586cf902f7644f4bcd98205d84fb897
      created: 1758235337837
      modified: 1758235337837
      sortKey: -1758235337837
      description: ""
    children:
      - url: "{{ base_url }}/api/v1/users"
        name: ➕ Criar Usuário
        meta:
          id: req_0ac4a5152f6543c4a30808dff67b6e85
          created: 1758235337835
          modified: 1758235337835
          isPrivate: false
          description: Criar um novo usuário no sistema (requer autenticação JWT)
          sortKey: -1758235337835
        method: POST
        body:
          mimeType: application/json
          text: |-
            {
              "name": "João Silva",
              "email": "joao.silva@empresa.com",
              "username": "joao.silva",
              "password": "senhaSegura123",
              "tenant_id": "uuid-do-tenant-aqui"
            }
        headers:
          - name: Authorization
            value: Bearer {{ auth_token }}
          - name: Content-Type
            value: application/json
          - name: Origin
            value: "{{ origin }}"
          - name: User-Agent
            value: "{{ user_agent }}"
          - name: Accept
            value: application/json
        settings:
          renderRequestBody: true
          encodeUrl: true
          followRedirects: global
          cookies:
            send: true
            store: true
          rebuildPath: true
  - name: 🏢 Tenants
    meta:
      id: fld_2a157e26b4ee42b5a32947a81e9c44ea
      created: 1758235337837
      modified: 1758235337837
      sortKey: -1758235337837
      description: ""
    children:
      - url: "{{ base_url }}/api/v1/tenants"
        name: 🏢 Criar Tenant
        meta:
          id: req_c84cb78802f04b92a96e9114f93755ed
          created: 1758235337836
          modified: 1758235337836
          isPrivate: false
          description: Criar um novo tenant/empresa no sistema (requer autenticação JWT)
          sortKey: -1758235337836
        method: POST
        body:
          mimeType: application/json
          text: |-
            {
              "name": "Empresa Exemplo Ltda",
              "email": "contato@empresa.com",
              "cpf_cnpj": "12345678000123",
              "type": "JURIDICA",
              "status": "ATIVO",
              "phone": "11999999999",
              "cell_phone": "11988888888",
              "street": "Rua das Flores",
              "number": "123",
              "complement": "Sala 456",
              "neighborhood": "Centro",
              "city": "São Paulo",
              "state": "SP",
              "cep": "01234567"
            }
        headers:
          - name: Authorization
            value: Bearer {{ auth_token }}
          - name: Content-Type
            value: application/json
          - name: Origin
            value: "{{ origin }}"
          - name: User-Agent
            value: "{{ user_agent }}"
          - name: Accept
            value: application/json
        settings:
          renderRequestBody: true
          encodeUrl: true
          followRedirects: global
          cookies:
            send: true
            store: true
          rebuildPath: true
cookieJar:
  name: Default Jar
  meta:
    id: jar_182cf93b4db5bfb942791411c7d316e1ef60003c
    created: 1758235337838
    modified: 1758235337838
environments:
  name: Base Environment
  meta:
    id: env_182cf93b4db5bfb942791411c7d316e1ef60003c
    created: 1758235337838
    modified: 1758235431430
    isPrivate: false
  data:
    base_url: http://localhost:5001
    api_key: VER_NO_DB
    origin: localhost
    auth_token: TOKEN_JWT
    refresh_token: REFRESH_TOKEN
EOF

# Criar arquivo de instruções atualizado
cat > ./docs/insomnia_import_instructions.md << 'EOF'
# 🌙 Como Importar para Insomnia

## 📋 Pré-requisitos
- [Insomnia](https://insomnia.rest/download) instalado
- Arquivos da documentação gerados

## 🚀 Passos para Importação
- Arquivos da documentação gerados

## 🚀 Passos para Importação

### Método 1: Importação Direta (Recomendado)

1. **Abra o Insomnia**
2. **Clique em "Create" > "Import/Export"**
3. **Selecione "Import Data" > "From File"**
4. **Escolha o arquivo:** `docs/insomnia_collection.yaml`
5. **Clique em "Import"**
6. **Configure a base URL:** `http://localhost:5001`

### Método 2: Via URL (Quando servidor estiver rodando)

1. **Abra o Insomnia**
2. **Clique em "Create" > "Import/Export"**
3. **Selecione "Import Data" > "From URL"**
4. **Cole a URL:** `http://localhost:5001/docs/insomnia_collection.json`
5. **Clique em "Fetch and Import"**

### Método 3: Usar arquivo Swagger original

1. **Abra o Insomnia**
2. **Clique em "Create" > "Import/Export"**
3. **Selecione "Import Data" > "From File"**
4. **Escolha o arquivo:** `docs/swagger.json`
5. **Clique em "Import"**

## 📚 Endpoints Disponíveis

Após a importação, você terá acesso a todos os endpoints:

### 🔐 Autenticação
- `POST /api/v1/auth/login` - Login de usuário
- `POST /api/v1/auth/refresh` - Renovar token
- `GET /api/v1/auth-apikey/tenant-by-apikey` - Buscar tenant por API Key

### 👥 Usuários
- `GET /api/v1/users` - Listar usuários
- `GET /api/v1/users/{id}` - Buscar usuário por ID
- `POST /api/v1/users` - Criar usuário
- `PUT /api/v1/users/{id}` - Atualizar usuário
- `PATCH /api/v1/users/{id}` - Atualizar parcialmente
- `DELETE /api/v1/users/{id}` - Excluir usuário

### 🏢 Tenants
- `GET /api/v1/tenants` - Listar tenants
- `GET /api/v1/tenants/{id}` - Buscar tenant por ID
- `POST /api/v1/tenants` - Criar tenant
- `PUT /api/v1/tenants/{id}` - Atualizar tenant
- `PATCH /api/v1/tenants/{id}` - Atualizar parcialmente
- `DELETE /api/v1/tenants/{id}` - Excluir tenant

## 💡 Dicas de Configuração

### Autenticação JWT
Para endpoints que requerem autenticação:
1. Faça login em `POST /api/v1/auth/login`
2. Copie o token da resposta
3. No Insomnia, vá em **Auth** > **Bearer Token**
4. Cole o token no campo **Token**

### Autenticação API Key
Para endpoints que usam X-API-Key:
1. Vá em **Auth** > **API Key Auth**
2. Configure:
   - **Key**: `X-API-Key`
   - **Value**: Sua API Key
   - **Add to**: `Header`

### Ambientes Disponíveis
A collection inclui 4 ambientes pré-configurados:

#### 🔧 Base Environment
- **Pré-configurado** com tokens JWT válidos e API key
- Ideal para **teste rápido** de endpoints autenticados
- **Tokens já incluídos** - não precisa fazer login
- **Variáveis completas**: base_url, api_key, auth_token, refresh_token, tenant_id

#### 💻 Local Development
- Ambiente completo para desenvolvimento local
- Inclui credenciais do master user
- Headers completos (Origin, User-Agent, etc.)

#### 🚀 Development Server
- Para testar em servidor de desenvolvimento
- Configure URL e credenciais conforme necessário

#### 🏭 Production
- Ambiente de produção
- Configure com credenciais reais

### Como Usar os Ambientes
1. No Insomnia, clique no dropdown de ambientes (canto superior direito)
2. Selecione o ambiente desejado
3. As variáveis serão automaticamente substituídas nas requests
4. Para o **Base Environment**, os tokens já estão configurados e funcionais

### Variáveis Disponíveis
Cada ambiente inclui estas variáveis:
- `base_url`: URL base da API (http://localhost:5001)
- `api_key`: Chave de API para autenticação X-API-Key
- `auth_token`: Token JWT de acesso (já configurado no Base Environment)
- `refresh_token`: Token JWT de refresh (já configurado no Base Environment)
- `origin`: Header Origin para CORS
- `user_agent`: Header User-Agent
- `master_email`: Email do usuário master (master@domain.local)
- `master_password`: Senha do usuário master (master123)
- `tenant_id`: ID do tenant atual (1c8c34eb-18e7-4b51-80bb-fc53ca1cc6e9)

## 🔗 Links Úteis
- [Documentação Swagger](http://localhost:5001/swagger/)
- [Coleção Postman](/docs/postman_collection.json)
- [Arquivo Insomnia](/docs/insomnia_collection.json)
- [Documentação JSON](/docs/swagger.json)
- [Documentação YAML](/docs/swagger.yaml)

---
*Gerado automaticamente pelo Go Base API*
EOF

echo "✅ Arquivos criados com sucesso!"
echo ""
echo "📁 Arquivos preparados:"
echo "   - docs/insomnia_collection.yaml (arquivo principal para importação)"
echo "   - docs/swagger.json (formato OpenAPI alternativo)"
echo "   - docs/swagger.yaml (formato YAML)"
echo "   - docs/insomnia_import_instructions.md (instruções completas)"
echo ""
echo "🌐 Próximos passos:"
echo "1. Abra o Insomnia"
echo "2. Importe o arquivo docs/insomnia_collection.yaml"
echo "3. Configure a base URL: http://localhost:5001"
echo "4. Comece a testar os endpoints!"
echo ""
echo "📖 Para mais detalhes, consulte: docs/insomnia_import_instructions.md"
echo ""
echo "🔗 Relacionado:"
echo "   - Documentação Swagger: make docs"
echo "   - Coleção Postman: make postman"
