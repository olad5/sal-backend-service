package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/olad5/sal-backend-service/internal/domain"
	"github.com/olad5/sal-backend-service/internal/infra"
)

var ErrMemoryStoreAccess = errors.New("error accessing memory store")

type MemoryProductRepository struct {
	products          map[uuid.UUID]domain.Product
	merchantsProducts map[uuid.UUID][]domain.Product
	lock              sync.RWMutex
}

func NewMemoryProductRepo() (*MemoryProductRepository, error) {
	return &MemoryProductRepository{
		map[uuid.UUID]domain.Product{},
		map[uuid.UUID][]domain.Product{},
		sync.RWMutex{},
	}, nil
}

func (m *MemoryProductRepository) CreateProduct(ctx context.Context, product domain.Product) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.products == nil {
		return ErrMemoryStoreAccess
	}
	if m.merchantsProducts == nil {
		return ErrMemoryStoreAccess
	}
	m.products[product.SKUID] = product
	m.merchantsProducts[product.MerchantId] = append(m.merchantsProducts[product.MerchantId], product)
	return nil
}

func (m *MemoryProductRepository) GetProductBySkuId(ctx context.Context, skuId uuid.UUID) (domain.Product, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.products == nil {
		return domain.Product{}, ErrMemoryStoreAccess
	}
	if m.merchantsProducts == nil {
		return domain.Product{}, ErrMemoryStoreAccess
	}

	existingProduct, err := m.getProductFromProductsStore(skuId)
	if err != nil {
		return domain.Product{}, err
	}
	return existingProduct, nil
}

func (m *MemoryProductRepository) GetProductsByMerchantId(ctx context.Context, merchantId uuid.UUID) ([]domain.Product, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.products == nil {
		return []domain.Product{}, ErrMemoryStoreAccess
	}
	if m.merchantsProducts == nil {
		return []domain.Product{}, ErrMemoryStoreAccess
	}
	if items, ok := m.merchantsProducts[merchantId]; ok {
		return items, nil
	}
	return []domain.Product{}, nil
}

func (m *MemoryProductRepository) UpdateProductByProductId(ctx context.Context, updatedProduct domain.Product) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.products == nil {
		return ErrMemoryStoreAccess
	}
	if m.merchantsProducts == nil {
		return ErrMemoryStoreAccess
	}

	m.products[updatedProduct.SKUID] = updatedProduct
	merchantProducts := m.merchantsProducts[updatedProduct.MerchantId]
	index, err := m.getIndexOfProduct(merchantProducts, updatedProduct.SKUID)
	if err != nil {
		return err
	}
	merchantProducts[index] = updatedProduct
	return nil
}

func (m *MemoryProductRepository) DeleteProductBySkuId(ctx context.Context, skuId uuid.UUID) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.products == nil {
		return ErrMemoryStoreAccess
	}
	if m.merchantsProducts == nil {
		return ErrMemoryStoreAccess
	}

	existingProduct, err := m.getProductFromProductsStore(skuId)
	if err != nil {
		return err
	}

	merchantProducts := m.merchantsProducts[existingProduct.MerchantId]
	index, err := m.getIndexOfProduct(merchantProducts, skuId)
	if err != nil {
		return err
	}
	if err := removeElement(merchantProducts, index); err != nil {
		return err
	}
	delete(m.products, skuId)

	return nil
}

func removeElement(slice []domain.Product, index int) error {
	if index < 0 || index >= len(slice) {
		return errors.New("out of bounds")
	}

	return nil
}

func (m *MemoryProductRepository) getProductFromProductsStore(skuId uuid.UUID) (domain.Product, error) {
	if existingProduct, ok := m.products[skuId]; ok {
		return existingProduct, nil
	}
	return domain.Product{}, infra.ErrProductNotFound
}

func (m *MemoryProductRepository) getIndexOfProduct(products []domain.Product, skuId uuid.UUID) (int, error) {
	for index, product := range products {
		if product.SKUID == skuId {
			return index, nil
		}
	}
	return 0, infra.ErrProductNotFound
}
