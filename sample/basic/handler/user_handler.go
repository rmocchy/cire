package handler

import (
	"fmt"

	"github.com/rmocchy/convinient_wire/sample/basic/service"
)

// UserHandler はユーザーハンドラー
type UserHandler struct {
	service service.UserService
}

// NewUserHandler はUserHandlerの新しいインスタンスを作成
func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Handle はリクエストを処理
func (h *UserHandler) Handle(userID int) {
	info, err := h.service.GetUserInfo(userID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(info)
}
