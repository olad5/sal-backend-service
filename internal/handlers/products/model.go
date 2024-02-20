package handlers

import (
	"time"

	"github.com/olad5/sal-backend-service/internal/domain"
)

type ProductDTO struct {
	SKUID       string     `json:"sku_id"`
	MerchantId  string     `json:"merchant_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func ToProductDTO(product domain.Product) ProductDTO {
	return ProductDTO{
		SKUID:       product.SKUID.String(),
		MerchantId:  product.MerchantId.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   &product.CreatedAt,
		UpdatedAt:   &product.UpdatedAt,
	}
}
