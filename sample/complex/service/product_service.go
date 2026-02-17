package service

import "github.com/rmocchy/cire/sample/complex/repository"

// ProductService は商品サービスのインターフェース
type ProductService interface {
	GetProduct(id int) string
}

// productServiceImpl はProductServiceの実装
type productServiceImpl struct {
	repo repository.ProductRepository
}

// NewProductService はProductServiceの新しいインスタンスを作成
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productServiceImpl{repo: repo}
}

func (s *productServiceImpl) GetProduct(id int) string {
	return s.repo.GetProductName(id)
}
