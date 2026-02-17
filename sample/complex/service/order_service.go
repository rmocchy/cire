package service

import (
	"fmt"

	"github.com/rmocchy/cire/sample/complex/repository"
)

// OrderService は注文サービスのインターフェース（並列依存の例）
type OrderService interface {
	CreateOrder(userID, productID int) string
}

// orderServiceImpl はOrderServiceの実装
// 複数のリポジトリに並列で依存している
type orderServiceImpl struct {
	userRepo    repository.UserRepository
	productRepo repository.ProductRepository
}

// NewOrderService はOrderServiceの新しいインスタンスを作成
// 並列依存: UserRepositoryとProductRepositoryの両方を必要とする
func NewOrderService(
	userRepo repository.UserRepository,
	productRepo repository.ProductRepository,
) OrderService {
	return &orderServiceImpl{
		userRepo:    userRepo,
		productRepo: productRepo,
	}
}

func (s *orderServiceImpl) CreateOrder(userID, productID int) string {
	userName := s.userRepo.GetUserName(userID)
	productName := s.productRepo.GetProductName(productID)
	return fmt.Sprintf("Order: %s purchased %s", userName, productName)
}
