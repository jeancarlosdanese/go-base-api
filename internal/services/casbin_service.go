// internal/services/casbin_service.go

// Package services inclui CasbinService, que configura e gerencia o controle de acesso baseado em políticas utilizando o Casbin.
// Esta implementação do CasbinService utiliza um modelo definido programaticamente e um adaptador GORM para interagir com o banco de dados.
// Com isso, eliminamos a necessidade de arquivos externos como 'model.conf', 'policy.csv' e 'config.go',
// e gerenciamos todas as políticas diretamente através da tabela 'casbin_rules' no banco de dados.
// Essa abordagem não só centraliza o gerenciamento das políticas de acesso como também reflete a tabela de permissões detalhada
// previamente representada no arquivo 'internal/config/casbin/policy.xlsx', garantindo que todas as políticas sejam aplicadas dinamicamente e possam ser modificadas em tempo real.

package services

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type CasbinService struct {
	enforcer *casbin.Enforcer
}

func NewCasbinService(db *gorm.DB) (*CasbinService, error) {
	// Configuração do modelo Casbin embutida diretamente no código.
	m, err := model.NewModelFromString(`
		[request_definition]
		r = sub, obj, act
		[policy_definition]
		p = sub, obj, act
		[role_definition]
		g = _, _
		[policy_effect]
		e = some(where (p.eft == allow))
		[matchers]
		m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act
	`)
	if err != nil {
		log.Printf("Erro ao carregar o modelo Casbin: %v", err)
		return nil, err
	}

	// Adaptador GORM para Casbin usando a tabela especificada.
	gormadapter.TurnOffAutoMigrate(db)
	a, err := gormadapter.NewAdapterByDBUseTableName(db, "casbin", "rules_view")
	if err != nil {
		log.Printf("Erro ao criar o adaptador GORM para Casbin: %v", err)
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(m, a)
	if err != nil {
		log.Printf("Erro ao criar o enforcer Casbin: %v", err)
		return nil, err
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		log.Printf("Erro ao carregar as políticas: %v", err)
		return nil, err
	}

	return &CasbinService{enforcer: enforcer}, nil
}

func (cs *CasbinService) CheckPermission(sub, obj, act string) bool {
	ok, err := cs.enforcer.Enforce(sub, obj, act)
	if err != nil {
		log.Printf("Erro ao verificar permissão: %v", err)
		return false
	}
	return ok
}
