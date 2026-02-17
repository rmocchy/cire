//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rmocchy/convinient_wire/sample/complex/handler"
	"github.com/rmocchy/convinient_wire/sample/complex/repository"
	"github.com/rmocchy/convinient_wire/sample/complex/service"
)

// UserAppSet はUserHandlerを持つルート構造体
type UserAppSet struct {
	UserHandler *handler.UserHandler
}

// OrderAppSet はProductHandlerとOrderHandlerを持つルート構造体（並列依存の例）
type OrderAppSet struct {
	ProductHandler *handler.ProductHandler
	OrderHandler   *handler.OrderHandler
}

// InitializeUserApp はUserHandlerの依存関係を解決して初期化
func InitializeUserApp() (*UserAppSet, error) {
	wire.Build(
		// Repository層
		repository.NewUserRepository,

		// Service層
		service.NewUserService,

		// Handler層
		handler.NewUserHandler,

		// ルート構造体
		wire.Struct(new(UserAppSet), "*"),
	)
	return nil, nil
}

// InitializeOrderApp はProductHandlerとOrderHandlerの依存関係を解決して初期化
// - OrderServiceは2つのリポジトリに並列で依存
func InitializeOrderApp() (*OrderAppSet, error) {
	wire.Build(
		// Repository層
		repository.NewUserRepository,
		repository.NewProductRepository,

		// Service層
		service.NewProductService,
		service.NewOrderService, // UserRepository + ProductRepositoryに並列依存

		// Handler層
		handler.NewProductHandler,
		handler.NewOrderHandler,

		// ルート構造体
		wire.Struct(new(OrderAppSet), "*"),
	)
	return nil, nil
}
