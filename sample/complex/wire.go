//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rmocchy/convinient_wire/sample/complex/handler"
	"github.com/rmocchy/convinient_wire/sample/complex/repository"
	"github.com/rmocchy/convinient_wire/sample/complex/service"
)

// UserAppSet is the Wire provider set for UserApp
var UserAppSet = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
	handler.NewUserHandler,
	wire.Struct(new(UserApp), "*"),
)

// InitializeUserApp initializes UserApp with all dependencies
func InitializeUserApp() (*UserApp, error) {
	wire.Build(UserAppSet)
	return nil, nil
}

// OrderAppSet is the Wire provider set for OrderApp
var OrderAppSet = wire.NewSet(
	repository.NewProductRepository,
	service.NewProductService,
	handler.NewProductHandler,
	repository.NewUserRepository,
	service.NewOrderService,
	handler.NewOrderHandler,
	wire.Struct(new(OrderApp), "*"),
)

// InitializeOrderApp initializes OrderApp with all dependencies
func InitializeOrderApp() (*OrderApp, error) {
	wire.Build(OrderAppSet)
	return nil, nil
}
