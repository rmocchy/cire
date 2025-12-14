package repository

import "fmt"

// Config はリポジトリの設定を表す構造体
type Config struct {
	DSN         string // データソース名
	MaxPoolSize int    // 最大接続プール数
}

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
	config *Config
}

// NewConfig はConfigの新しいインスタンスを作成
func NewConfig() *Config {
	return &Config{
		DSN:         "user:password@tcp(localhost:3306)/mydb",
		MaxPoolSize: 10,
	}
}

// NewUserRepository はUserRepositoryの新しいインスタンスを作成
func NewUserRepository(config *Config) (UserRepository, error) {
	return &userRepositoryImpl{
		config: config,
	}, nil
}

func (r *userRepositoryImpl) FindByID(id int) (*User, error) {
	// 実際にはDBからデータを取得
	// ここではダミーデータを返す
	return &User{
		ID:   id,
		Name: fmt.Sprintf("User%d", id),
	}, nil
}
