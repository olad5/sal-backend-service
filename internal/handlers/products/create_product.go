package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/olad5/sal-backend-service/internal/usecases/products"
	appErrors "github.com/olad5/sal-backend-service/pkg/errors"

	"github.com/olad5/sal-backend-service/pkg/utils"
)

func (p ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Body == nil {
		utils.ErrorResponse(w, appErrors.ErrMissingBody, http.StatusBadRequest)
		return
	}
	type requestDTO struct {
		SKUID       string  `json:"sku_id"`
		MerchantId  string  `json:"merchant_id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	var request requestDTO
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.ErrorResponse(w, appErrors.ErrInvalidJson, http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		utils.ErrorResponse(w, "name required", http.StatusBadRequest)
		return
	}
	if request.Description == "" {
		utils.ErrorResponse(w, "description required", http.StatusBadRequest)
		return
	}
	if request.Price < 0 {
		utils.ErrorResponse(w, "price cannot be less than zero", http.StatusBadRequest)
		return
	}

	merchantId, err := uuid.Parse(request.MerchantId)
	if err != nil {
		utils.ErrorResponse(w, appErrors.ErrInvalidID.Error(), http.StatusBadRequest)
		return
	}
	skuId, err := uuid.Parse(request.SKUID)
	if err != nil {
		utils.ErrorResponse(w, appErrors.ErrInvalidID.Error(), http.StatusBadRequest)
		return
	}

	newProduct, err := p.productService.CreateProduct(ctx, merchantId, skuId, request.Name, request.Description, request.Price)
	if err != nil {
		switch {
		case errors.Is(err, products.ErrProductAlreadyExists):
			utils.ErrorResponse(w, products.ErrProductAlreadyExists.Error(), http.StatusBadRequest)
			return
		default:
			utils.ErrorResponse(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError)
			return
		}
	}

	utils.SuccessResponse(w, "product created successfully", ToProductDTO(newProduct))
}
