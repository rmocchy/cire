package main

import (
	"github.com/rmocchy/cire/sample/complex/handler"
)

// UserApp はユーザー関係のハンドラーを持つルート構造体
type UserApp struct {
	UserHandler *handler.UserHandler
}

// OrderApp は注文関係のハンドラーを持つルート構造体
type OrderApp struct {
	ProductHandler *handler.ProductHandler
	OrderHandler   *handler.OrderHandler
}
