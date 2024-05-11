// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package config

import (
	"github.com/jeancarlosdanese/go-base-api/internal/app"
)

// Injectors from wire.go:

func InitializeServicesContainer() (*app.ServicesContainer, error) {
	servicesContainer, err := app.NewServicesContainer()
	if err != nil {
		return nil, err
	}
	return servicesContainer, nil
}
