package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/olad5/sal-backend-service/internal/infra"
	"github.com/olad5/sal-backend-service/internal/usecases/products"
	appErrors "github.com/olad5/sal-backend-service/pkg/errors"

	"github.com/olad5/sal-backend-service/pkg/utils"
)

func (p ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
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
		MerchantId string `json:"merchant_id"`
	}

	var request requestDTO
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.ErrorResponse(w, appErrors.ErrInvalidJson, http.StatusBadRequest)
		return
	}

	merchantId, err := uuid.Parse(request.MerchantId)
	if err != nil {
		utils.ErrorResponse(w, appErrors.ErrInvalidID.Error(), http.StatusBadRequest)
		return
	}

	err = p.productService.DeleteProduct(ctx, merchantId, skuId)
	if err != nil {
		switch {
		case errors.Is(err, infra.ErrProductNotFound):
			utils.ErrorResponse(w, infra.ErrProductNotFound.Error(), http.StatusNotFound)
			return
		case errors.Is(err, products.ErrUserNotAuthorized):
			utils.ErrorResponse(w, appErrors.ErrUnauthorized, http.StatusNotFound)
			return
		default:
			utils.ErrorResponse(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError)
			return
		}
	}

	utils.SuccessResponse(w, "product deleted successfully", nil)
}
