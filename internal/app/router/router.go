package router

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	handlers "github.com/olad5/sal-backend-service/internal/handlers/products"
	"github.com/olad5/sal-backend-service/internal/infra/memory"
	"github.com/olad5/sal-backend-service/internal/usecases/products"
)

func NewHttpRouter(ctx context.Context) http.Handler {
	productRepo, err := memory.NewMemoryProductRepo()
	if err != nil {
		log.Fatal("Error Initializing Product Repo", err)
	}

	productService, err := products.NewProductService(productRepo)
	if err != nil {
		log.Fatal("Error Initializing ProductService")
	}

	productHandler, err := handlers.NewProductHandler(*productService)
	if err != nil {
		log.Fatal("failed to create the Product handler: ", err)
	}
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "SAL Backend Service is live\n")
	})

	router.Route("/api", func(r chi.Router) {
		r.Use(
			middleware.AllowContentType("application/json"),
			middleware.SetHeader("Content-Type", "application/json"),
		)
		r.Post("/products", productHandler.CreateProduct)
		r.Patch("/products/{sku_id}", productHandler.EditProduct)
		r.Delete("/products/{sku_id}", productHandler.DeleteProduct)
		r.Get("/merchants/{merchant_id}/products", productHandler.FetchMerchantProducts)
	})
	return router
}
