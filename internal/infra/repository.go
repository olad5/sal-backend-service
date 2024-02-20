package infra

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/olad5/sal-backend-service/internal/domain"
)

var ErrProductNotFound = errors.New("product not found")

type ProductRepository interface {
	CreateProduct(ctx context.Context, product domain.Product) error
	GetProductBySkuId(ctx context.Context, skuId uuid.UUID) (domain.Product, error)
	GetProductsByMerchantId(ctx context.Context, merchantId uuid.UUID) ([]domain.Product, error)
	UpdateProductByProductId(ctx context.Context, product domain.Product) error
	DeleteProductBySkuId(ctx context.Context, skuId uuid.UUID) error
}
