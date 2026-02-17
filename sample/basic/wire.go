//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rmocchy/convinient_wire/sample/basic/handler"
	"github.com/rmocchy/convinient_wire/sample/basic/repository"
	"github.com/rmocchy/convinient_wire/sample/basic/service"
)

// ControllerSet は依存関係の解析対象となるルート構造体
type ControllerSet struct {
	handler *handler.UserHandler
}

// InitializeUserHandler は全ての依存関係を解決してUserHandlerを初期化
func InitializeUserHandler() (*ControllerSet, error) {
	wire.Build(
		repository.NewConfig,
		repository.NewUserRepository,
		service.NewUserService,
		handler.NewUserHandler,
		wire.Struct(new(ControllerSet), "*"),
	)
	return nil, nil
}
