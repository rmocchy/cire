package service

import (
	"fmt"

	"github.com/rmocchy/cire/sample/basic/repository"
)

// UserService はユーザーサービスのインターフェース
type UserService interface {
	GetUserInfo(id int) (string, error)
}

// userServiceImpl はUserServiceの実装
type userServiceImpl struct {
	repo repository.UserRepository
}

// NewUserService はUserServiceの新しいインスタンスを作成
func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{
		repo: repo,
	}
}

func (s *userServiceImpl) GetUserInfo(id int) (string, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("User Info: ID=%d, Name=%s", user.ID, user.Name), nil
}
