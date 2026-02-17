//go:build wireinject
// +build wireinject

package main

import (
	"github.com/rmocchy/convinient_wire/sample/complex/handler"
)

// UserApp はUserHandlerを持つルート構造体
type UserApp struct {
	UserHandler *handler.UserHandler
}

// OrderApp はProductHandlerとOrderHandlerを持つルート構造体（並列依存の例）
type OrderApp struct {
	ProductHandler *handler.ProductHandler
	OrderHandler   *handler.OrderHandler
}
