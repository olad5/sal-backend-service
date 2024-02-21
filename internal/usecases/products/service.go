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

func (p *ProductService) CreateProduct(ctx context.Context, merchantId, skuId uuid.UUID, name, description string, price float64) (domain.Product, error) {
	existingProduct, err := p.productRepo.GetProductBySkuId(ctx, skuId)
	if err == nil && existingProduct.SKUID == skuId {
		return domain.Product{}, ErrProductAlreadyExists
	}

	newProduct := domain.Product{
		SKUID:       skuId,
		MerchantId:  merchantId,
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = p.productRepo.CreateProduct(ctx, newProduct)
	if err != nil {
		return domain.Product{}, err
	}
	return newProduct, nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, merchantId, skuId uuid.UUID, name, description string, price float64) (domain.Product, error) {
	existingProduct, err := p.productRepo.GetProductBySkuId(ctx, skuId)
	if err != nil {
		return domain.Product{}, err
	}
	var updatedName string
	if name != "" {
		updatedName = name
	} else {
		updatedName = existingProduct.Name
	}

	var updatedDescription string
	if name != "" {
		updatedDescription = description
	} else {
		updatedDescription = existingProduct.Description
	}

	var updatedPrice float64
	if name != "" {
		updatedPrice = price
	} else {
		updatedPrice = existingProduct.Price
	}

	updatedProduct := domain.Product{
		SKUID:       existingProduct.SKUID,
		Name:        updatedName,
		Description: updatedDescription,
		Price:       updatedPrice,
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
	_, err := p.productRepo.GetProductBySkuId(ctx, productId)
	if err != nil {
		return err
	}

	err = p.productRepo.DeleteProductBySkuId(ctx, productId)
	if err != nil {
		return err
	}

	return nil
}
