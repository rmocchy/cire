package handler

import "github.com/rmocchy/convinient_wire/sample/complex/service"

// UserHandler はユーザーハンドラー
type UserHandler struct {
	service service.UserService
}

// NewUserHandler はUserHandlerの新しいインスタンスを作成
func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Handle はリクエストを処理
func (h *UserHandler) Handle(userID int) string {
	return h.service.GetUser(userID)
}
