package handler

import "github.com/rmocchy/cire/sample/complex/service"

// OrderHandler は注文ハンドラー
type OrderHandler struct {
	service service.OrderService
}

// NewOrderHandler はOrderHandlerの新しいインスタンスを作成
func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// Handle はリクエストを処理
func (h *OrderHandler) Handle(userID, productID int) string {
	return h.service.CreateOrder(userID, productID)
}
