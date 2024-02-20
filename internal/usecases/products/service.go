package products

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/olad5/sal-backend-service/internal/domain"
	"github.com/olad5/sal-backend-service/internal/infra"

	"github.com/google/uuid"
)

type ProductService struct {
	productRepo infra.ProductRepository
}

var ErrProductAlreadyExists = errors.New("product already exists")

func NewProductService(productRepo infra.ProductRepository) (*ProductService, error) {
	if productRepo == nil {
		return &ProductService{}, fmt.Errorf("ProductService failed to initialize, productRepo is nil")
	}
	return &ProductService{productRepo}, nil
}

func (p *ProductService) CreateProduct(ctx context.Context, merchantId, skuId uuid.UUID, name, description string, price int64) (domain.Product, error) {
	existingProduct, err := p.productRepo.GetProductByProductId(ctx, skuId)
	if err == nil && existingProduct.SKUID == skuId {
		return domain.Product{}, ErrProductAlreadyExists
	}

	newProduct := domain.Product{
		SKUID:       skuId,
		Name:        name,
		Description: description,
		Price:       price,
		MerchantId:  merchantId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = p.productRepo.CreateProduct(ctx, newProduct)
	if err != nil {
		return domain.Product{}, err
	}
	return newProduct, nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, merchantId, skuId uuid.UUID, name, description string, price int64) (domain.Product, error) {
	existingProduct, err := p.productRepo.GetProductByProductId(ctx, skuId)
	if err != nil {
		return domain.Product{}, err
	}

	updatedProduct := domain.Product{
		SKUID:       existingProduct.SKUID,
		Name:        "",
		Description: "",
		Price:       0,
		MerchantId:  existingProduct.MerchantId,
		CreatedAt:   existingProduct.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	err = p.productRepo.UpdateProductByProductId(ctx, updatedProduct)
	if err != nil {
		return domain.Product{}, err
	}
	return updatedProduct, nil
}

func (p *ProductService) GetProductsByMerchantId(ctx context.Context, merchantId uuid.UUID) ([]domain.Product, error) {
	products, err := p.productRepo.GetProductsByMerchantId(ctx, merchantId)
	if err != nil {
		return []domain.Product{}, err
	}

	return products, nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, productId uuid.UUID) error {
	_, err := p.productRepo.GetProductByProductId(ctx, productId)
	if err != nil {
		return err
	}

	err = p.productRepo.DeleteProductByProductId(ctx, productId)
	if err != nil {
		return err
	}

	return nil
}
