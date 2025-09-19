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
