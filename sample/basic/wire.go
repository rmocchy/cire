//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rmocchy/convinient_wire/sample/basic-di/handler"
	"github.com/rmocchy/convinient_wire/sample/basic-di/repository"
	"github.com/rmocchy/convinient_wire/sample/basic-di/service"
)

type ControllerSet struct {
	handler *handler.UserHandler
}

// InitializeUserHandler は全ての依存関係を解決してUserHandlerを初期化
func InitializeUserHandler() (*ControllerSet, error) {
	wire.Build(
		repository.NewUserRepository,
		service.NewUserService,
		handler.NewUserHandler,
		wire.Struct(new(ControllerSet), "*"),
	)
	return nil, nil
}
