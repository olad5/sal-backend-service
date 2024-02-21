package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/olad5/sal-backend-service/internal/domain"
	"github.com/olad5/sal-backend-service/internal/infra"
)

type MemoryProductRepository struct {
	store map[uuid.UUID]domain.Product
	lock  sync.RWMutex
}

func NewMemoryProductRepo() (*MemoryProductRepository, error) {
	return &MemoryProductRepository{
		map[uuid.UUID]domain.Product{},
		sync.RWMutex{},
	}, nil
}

func (m *MemoryProductRepository) CreateProduct(ctx context.Context, product domain.Product) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.store[product.SKUID] = product
	return nil
}

func (m *MemoryProductRepository) GetProductBySkuId(ctx context.Context, skuId uuid.UUID) (domain.Product, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if existingProduct, ok := m.store[skuId]; ok {
		return existingProduct, nil
	}
	return domain.Product{}, infra.ErrProductNotFound
}

func (m *MemoryProductRepository) GetProductsByMerchantId(ctx context.Context, merchantId uuid.UUID) ([]domain.Product, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	results := []domain.Product{}
	for _, product := range m.store {
		if product.MerchantId == merchantId {
			results = append(results, product)
		}
	}
	return results, nil
}

func (m *MemoryProductRepository) UpdateProductByProductId(ctx context.Context, updatedProduct domain.Product) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.store[updatedProduct.SKUID] = updatedProduct
	return nil
}

func (m *MemoryProductRepository) DeleteProductBySkuId(ctx context.Context, skuId uuid.UUID) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.store, skuId)

	return nil
}
