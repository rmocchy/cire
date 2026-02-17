//go:build cire
// +build cire

package main

import (
	"github.com/rmocchy/convinient_wire/sample/complex/handler"
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
