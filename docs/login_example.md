# 📝 Exemplos de Login

## 🚀 Como Fazer Login

### Usando cURL

```bash
# ✅ Login básico - FORMA CORRETA
curl -X POST http://localhost:5001/api/v1/auth/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=master@domain.local&password=master123"

# ✅ Com Origin explícito (caso necessário)
curl -X POST http://localhost:5001/api/v1/auth/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -H "Origin: http://localhost" \
  -d "email=master@domain.local&password=master123"
```

### Usando Postman

1. **Método:** POST
2. **URL:** `http://localhost:5001/api/v1/auth/login`
3. **Headers:**
   - `Content-Type: application/x-www-form-urlencoded`
4. **Body (form-data):**
   - `email`: `master@domain.local`
   - `password`: `master123`

### Resposta de Sucesso

```json
{
  "type": "Bearer",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid-do-usuario",
    "name": "Master User",
    "username": "master",
    "email": "master@domain.local"
  },
  "roles": ["master"],
  "policies": ["GET|POST|PUT|PATCH|DELETE"]
}
```

## 🔧 Configuração do Ambiente

### Headers Necessários

Para requisições autenticadas, use:
```bash
Authorization: Bearer SEU_TOKEN_AQUI
```

### Refresh Token

Para renovar o token:
```bash
curl -X POST http://localhost:5001/api/v1/auth/refresh \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "refreshToken=SEU_REFRESH_TOKEN_AQUI"
```

## ⚠️ Possíveis Erros

### Erro 400 - Bad Request
```json
{
  "error": "Parâmetros de entrada inválidos",
  "details": "Key: 'LoginForm.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'LoginForm.Password' Error:Field validation for 'Password' failed on the 'required' tag",
  "hint": "Verifique se os campos 'email' e 'password' estão presentes no form-data"
}
```

**Causa:** Content-Type incorreto ou campos não enviados como form-data.

**Solução:** Sempre use `Content-Type: application/x-www-form-urlencoded` e envie dados como `key=value&key2=value2`.

### Erro 401 - Unauthorized
```json
{
  "error": "usuário ou origem não encontrado",
  "hint": "Verifique se o email e senha estão corretos"
}
```

### Erro 400 - Origin não fornecida
```json
{
  "error": "Origem não fornecida",
  "hint": "Certifique-se de que o cliente HTTP está enviando o header Origin ou configure localhost como padrão"
}
```

## 🐛 Debug

### Verificar se o servidor está rodando
```bash
curl http://localhost:5001/health
```

### Verificar se o usuário existe
```bash
# Verifique se o usuário master foi criado na migração
# Email: master@domain.local
# Senha: master123 (ou a senha configurada em $MASTER_PASSWORD_HASH)
```

### Testar com Origin explícito
```bash
curl -X POST http://localhost:5001/api/v1/auth/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -H "Origin: http://localhost" \
  -d "email=master@domain.local&password=master123"
```

### 🔧 Correção Específica do Erro Mostrado

**Se você recebeu exatamente este erro:**
```json
{
	"details": "Key: 'LoginForm.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'LoginForm.Password' Error:Field validation for 'Password' failed on the 'required' tag",
	"error": "Parâmetros de entrada inválidos",
	"hint": "Verifique se os campos 'email' e 'password' estão presentes no form-data"
}
```

**Use EXATAMENTE este comando:**
```bash
curl -X POST http://localhost:5001/api/v1/auth/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=master@domain.local&password=master123"
```

**Diferenças críticas:**
- ✅ `Content-Type: application/x-www-form-urlencoded` (obrigatório)
- ✅ Dados como `key=value&key=value2` (não JSON)
- ✅ Campos `email` e `password` presentes e não vazios

---
*Para mais exemplos, consulte a documentação Swagger em http://localhost:5001/swagger/*
