package service

import "github.com/rmocchy/convinient_wire/sample/complex/repository"

// UserService はユーザーサービスのインターフェース
type UserService interface {
	GetUser(id int) string
}

// userServiceImpl はUserServiceの実装
type userServiceImpl struct {
	repo repository.UserRepository
}

// NewUserService はUserServiceの新しいインスタンスを作成
func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) GetUser(id int) string {
	return s.repo.GetUserName(id)
}
