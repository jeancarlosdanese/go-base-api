// internal/config/wire.go

//go:build wireinject
// +build wireinject

package config

import (
	"github.com/google/wire"
	"github.com/jeancarlosdanese/go-base-api/internal/app"
)

func InitializeServicesContainer() (*app.ServicesContainer, error) {
	wire.Build(
		app.NewServicesContainer,
	)
	return nil, nil
}
