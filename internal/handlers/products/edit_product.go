package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/olad5/sal-backend-service/internal/infra"
	appErrors "github.com/olad5/sal-backend-service/pkg/errors"

	"github.com/olad5/sal-backend-service/pkg/utils"
)

func (p ProductHandler) EditProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "sku_id")
	if id == "" {
		utils.ErrorResponse(w, "sku_id required", http.StatusBadRequest)
		return
	}

	skuId, err := uuid.Parse(id)
	if err != nil {
		utils.ErrorResponse(w, appErrors.ErrInvalidID.Error(), http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		utils.ErrorResponse(w, appErrors.ErrMissingBody, http.StatusBadRequest)
		return
	}
	type requestDTO struct {
		MerchantId  string  `json:"merchant_id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	var request requestDTO
	err = json.NewDecoder(r.Body).Decode(&request)
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

	updatedProduct, err := p.productService.UpdateProduct(ctx, merchantId, skuId, request.Name, request.Description, request.Price)
	if err != nil {
		switch {
		case errors.Is(err, infra.ErrProductNotFound):
			utils.ErrorResponse(w, infra.ErrProductNotFound.Error(), http.StatusNotFound)
			return
		default:
			utils.ErrorResponse(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError)
			return
		}
	}

	utils.SuccessResponse(w, "product updated successfully", ToProductDTO(updatedProduct))
}
