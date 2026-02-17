package handler

import "github.com/rmocchy/convinient_wire/sample/complex/service"

// ProductHandler は商品ハンドラー
type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler はProductHandlerの新しいインスタンスを作成
func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// Handle はリクエストを処理
func (h *ProductHandler) Handle(productID int) string {
	return h.service.GetProduct(productID)
}
