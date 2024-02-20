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
	GetProductByProductId(ctx context.Context, productId uuid.UUID) (domain.Product, error)
	UpdateProductByProductId(ctx context.Context, product domain.Product) error
	DeleteProductByProductId(ctx context.Context, product domain.Product) error
}
