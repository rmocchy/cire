package service

import "fmt"

// UserService はユーザーサービスのインターフェース
type UserService interface {
	GetUserInfo(id int) (string, error)
}

type userServiceImpl struct{}

// NewUserService はUserServiceの新しいインスタンスを作成
func NewUserService() UserService {
	return &userServiceImpl{}
}

// NewAltUserService はUserServiceの別コンストラクタ
// NOTE: NewUserService と同じ型を返すため、cire は重複エラーを検出して失敗します
func NewAltUserService() UserService {
	return &userServiceImpl{}
}

func (s *userServiceImpl) GetUserInfo(id int) (string, error) {
	return fmt.Sprintf("UserInfo: id=%d", id), nil
}
