#!/bin/bash

# Script para atualizar a documentação Swagger
# Uso: ./scripts/update_docs.sh

echo "🚀 Atualizando documentação Swagger..."

# Verificar se o swag está instalado
if ! command -v swag &> /dev/null; then
    echo "❌ swag não encontrado. Instalando..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Mudar para o diretório do projeto se não estiver lá
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR" || exit 1

# Executar swag init da raiz do projeto (funciona melhor)
echo "📝 Gerando documentação..."
swag init -g ./cmd/go_api/main.go --output ./docs/ --parseDependency --parseInternal --parseDepth 2 --useStructName

if [ $? -eq 0 ]; then
    echo "✅ Documentação Swagger atualizada com sucesso!"
    echo "📁 Arquivos gerados:"
    ls -la ./docs/
    echo ""
    echo "🌐 Acesse a documentação em:"
    echo "   - Interface Web: http://localhost:5001/swagger/"
    echo "   - JSON: http://localhost:5001/swagger/doc.json"
    echo "   - YAML: http://localhost:5001/swagger/swagger.yaml"
else
    echo "❌ Erro ao gerar documentação!"
    exit 1
fi
