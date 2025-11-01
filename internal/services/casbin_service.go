// internal/services/casbin_service.go

// Package services inclui CasbinService, que configura e gerencia o controle de acesso baseado em políticas utilizando o Casbin.
// Esta implementação do CasbinService utiliza um modelo definido programaticamente e um adaptador GORM para interagir com o banco de dados.
// Com isso, eliminamos a necessidade de arquivos externos como 'model.conf', 'policy.csv' e 'config.go',
// e gerenciamos todas as políticas diretamente através da tabela 'casbin_rules' no banco de dados.
// Essa abordagem não só centraliza o gerenciamento das políticas de acesso como também reflete a tabela de permissões detalhada
// previamente representada no arquivo 'internal/config/casbin/policy.xlsx', garantindo que todas as políticas sejam aplicadas dinamicamente e possam ser modificadas em tempo real.

package services

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/jeancarlosdanese/go-base-api/internal/logging"
	"gorm.io/gorm"
)

type CasbinServiceInterface interface {
	CheckPermission(sub, obj, act string) bool
}
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
		[policy_effect]
		e = some(where (p.eft == allow))
		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
	`)
	if err != nil {
		logging.Logger.Error().
			Err(err).
			Str("operation", "load_casbin_model").
			Msg("Erro ao carregar o modelo Casbin")
		return nil, err
	}

	// Adaptador GORM para Casbin usando a tabela especificada.
	gormadapter.TurnOffAutoMigrate(db)
	a, err := gormadapter.NewAdapterByDBUseTableName(db, "casbin", "rules_view")
	if err != nil {
		logging.Logger.Error().
			Err(err).
			Str("operation", "create_casbin_adapter").
			Msg("Erro ao criar o adaptador GORM para Casbin")
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(m, a)
	if err != nil {
		logging.Logger.Error().
			Err(err).
			Str("operation", "create_casbin_enforcer").
			Msg("Erro ao criar o enforcer Casbin")
		return nil, err
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		logging.Logger.Error().
			Err(err).
			Str("operation", "load_casbin_policies").
			Msg("Erro ao carregar as políticas")
		return nil, err
	}

	logging.Logger.Info().
		Str("operation", "init_casbin").
		Msg("Casbin inicializado com sucesso")

	return &CasbinService{enforcer: enforcer}, nil
}

// CheckPermission logs and verifies permissions using Casbin enforcer
func (cs *CasbinService) CheckPermission(sub, obj, act string) bool {
	ok, err := cs.enforcer.Enforce(sub, obj, act)
	if err != nil {
		logging.Logger.Error().
			Err(err).
			Str("subject", sub).
			Str("object", obj).
			Str("action", act).
			Str("operation", "check_permission").
			Msg("Erro ao verificar permissão")
		return false
	}

	logging.Logger.Debug().
		Str("subject", sub).
		Str("object", obj).
		Str("action", act).
		Bool("allowed", ok).
		Str("operation", "check_permission").
		Msg("Permissão verificada")
	return ok
}
