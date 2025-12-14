package repository

import "fmt"

// User はユーザー情報を表す構造体
type User struct {
	ID   int
	Name string
}

// UserRepository はユーザーリポジトリのインターフェース
type UserRepository interface {
	FindByID(id int) (*User, error)
}

// userRepositoryImpl はUserRepositoryの実装
type userRepositoryImpl struct {
	// 実際にはDBコネクションなどを持つ
}

// NewUserRepository はUserRepositoryの新しいインスタンスを作成
func NewUserRepository() (UserRepository, error) {
	return &userRepositoryImpl{}, nil
}

func (r *userRepositoryImpl) FindByID(id int) (*User, error) {
	// 実際にはDBからデータを取得
	// ここではダミーデータを返す
	return &User{
		ID:   id,
		Name: fmt.Sprintf("User%d", id),
	}, nil
}
