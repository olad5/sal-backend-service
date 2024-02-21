package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/olad5/sal-backend-service/internal/infra"
	"github.com/olad5/sal-backend-service/internal/usecases/products"
	appErrors "github.com/olad5/sal-backend-service/pkg/errors"

	"github.com/olad5/sal-backend-service/pkg/utils"
)

func (p ProductHandler) FetchMerchantProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "merchant_id")
	if id == "" {
		utils.ErrorResponse(w, "merchant_id required", http.StatusBadRequest)
		return
	}

	merchantId, err := uuid.Parse(id)
	if err != nil {
		utils.ErrorResponse(w, appErrors.ErrInvalidID.Error(), http.StatusBadRequest)
		return
	}

	merchantProducts, err := p.productService.GetProductsByMerchantId(ctx, merchantId)
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

	utils.SuccessResponse(w, "products retrieved successfully", ToProductPagedDTO(merchantProducts))
}
