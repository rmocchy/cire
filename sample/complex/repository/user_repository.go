package repository

// UserRepository はユーザーリポジトリのインターフェース
type UserRepository interface {
	GetUserName(id int) string
}

// userRepositoryImpl はUserRepositoryの実装
type userRepositoryImpl struct{}

// NewUserRepository はUserRepositoryの新しいインスタンスを作成
func NewUserRepository() UserRepository {
	return &userRepositoryImpl{}
}

func (r *userRepositoryImpl) GetUserName(id int) string {
	return "User"
}
