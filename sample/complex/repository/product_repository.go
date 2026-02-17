package repository

// ProductRepository は商品リポジトリのインターフェース
type ProductRepository interface {
	GetProductName(id int) string
}

// productRepositoryImpl はProductRepositoryの実装
type productRepositoryImpl struct{}

// NewProductRepository はProductRepositoryの新しいインスタンスを作成
func NewProductRepository() ProductRepository {
	return &productRepositoryImpl{}
}

func (r *productRepositoryImpl) GetProductName(id int) string {
	return "Product"
}
