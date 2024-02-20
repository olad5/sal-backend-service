package handlers

import (
	"errors"

	"github.com/olad5/sal-backend-service/internal/usecases/products"
)

type ProductHandler struct {
	productService products.ProductService
}

func NewProductHandler(productService products.ProductService) (*ProductHandler, error) {
	if productService == (products.ProductService{}) {
		return nil, errors.New("product service cannot be empty")
	}

	return &ProductHandler{productService}, nil
}
