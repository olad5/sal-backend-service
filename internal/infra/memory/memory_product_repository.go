package memory

import (
	"context"

	"github.com/google/uuid"
	"github.com/olad5/sal-backend-service/internal/domain"
	"github.com/olad5/sal-backend-service/internal/infra"
)

type MemoryProductRepository struct {
	products []domain.Product
}

func NewMemoryProductRepo() (*MemoryProductRepository, error) {
	return &MemoryProductRepository{products: []domain.Product{}}, nil
}

func (m *MemoryProductRepository) CreateProduct(ctx context.Context, product domain.Product) error {
	m.products = append(m.products, product)
	return nil
}

func (m *MemoryProductRepository) GetProductByProductId(ctx context.Context, productId uuid.UUID) (domain.Product, error) {
	indexOfExistingProduct, err := m.getIndexOfProduct(productId)
	if err != nil {
		return domain.Product{}, err
	}
	return m.products[indexOfExistingProduct], nil
}

func (m *MemoryProductRepository) GetProductsByMerchantId(ctx context.Context, merchantId uuid.UUID) ([]domain.Product, error) {
	results := []domain.Product{}
	for _, product := range m.products {
		if product.MerchantId == merchantId {
			results = append(results, product)
		}
	}
	return results, nil
}

func (m *MemoryProductRepository) UpdateProductByProductId(ctx context.Context, updatedProduct domain.Product) error {
	indexOfExistingProduct, err := m.getIndexOfProduct(updatedProduct.SKUID)
	if err != nil {
		return err
	}
	m.products[indexOfExistingProduct] = updatedProduct
	return nil
}

func (m *MemoryProductRepository) DeleteProductByProductId(ctx context.Context, productId uuid.UUID) error {
	indexToRemove, err := m.getIndexOfProduct(productId)
	if err != nil {
		return err
	}

	existingProducts := m.products
	newSlice := []domain.Product{}
	m.products = append(newSlice, existingProducts[:indexToRemove]...)
	m.products = append(newSlice, existingProducts[indexToRemove:]...)

	return nil
}

func (m *MemoryProductRepository) getIndexOfProduct(productId uuid.UUID) (int, error) {
	for index, product := range m.products {
		if product.SKUID == productId {
			return index, nil
		}
	}
	return 0, infra.ErrProductNotFound
}
