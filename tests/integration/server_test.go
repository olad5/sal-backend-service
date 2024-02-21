//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/olad5/sal-backend-service/internal/app/router"
	"github.com/olad5/sal-backend-service/tests"
)

var r http.Handler

func TestMain(m *testing.M) {
	ctx := context.Background()
	r = router.NewHttpRouter(ctx)

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCreateProduct(t *testing.T) {
	route := "/api/products"
	type Product struct {
		SKUID       uuid.UUID `json:"sku_id"`
		MerchantId  uuid.UUID `json:"merchant_id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Price       float64   `json:"price"`
	}
	t.Run("test for invalid json request body",
		func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, route, nil)
			response := tests.ExecuteRequest(req, r)
			tests.AssertStatusCode(t, http.StatusBadRequest, response.Code)
		},
	)
	t.Run(`Given a merchant wants to create a new product with valid attributes,
    when they make a POST request to the create product endpoint with valid data,
    then the product should be successfully created in the database. `,
		func(t *testing.T) {
			merchantId := uuid.New()
			skuId := uuid.New()
			var price float64 = 30.00
			np := Product{
				MerchantId:  merchantId,
				SKUID:       skuId,
				Name:        "some-product-name",
				Description: "some-product-description",
				Price:       price,
			}
			requestBody, err := json.Marshal(&np)
			if err != nil {
				t.Fatal(err)
			}
			req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(requestBody))
			response := tests.ExecuteRequest(req, r)
			tests.AssertStatusCode(t, http.StatusOK, response.Code)
			data := tests.ParseResponse(t, response)["data"].(map[string]interface{})
			tests.AssertResponseMessage(t, data["merchant_id"].(string), merchantId.String())
			tests.AssertResponseMessage(t, data["sku_id"].(string), skuId.String())
			if responsePrice, ok := data["price"].(float64); ok {
				tests.AssertResponseMessage(t, strconv.FormatFloat(responsePrice, 'f', -1, 64), strconv.FormatFloat(price, 'f', -1, 64))
			}
		},
	)

	t.Run(`Given a merchant tries to create a product with a SKU ID that already exists,
        when they make a POST request to the create product endpoint,
        then the API should return a conflict error indicating that the SKU ID 
        already exists. `,
		func(t *testing.T) {
			merchantId := uuid.New()
			skuId := uuid.New()
			var price float64 = 10.00
			np := Product{
				MerchantId:  merchantId,
				SKUID:       skuId,
				Name:        "some-product-name",
				Description: "some-product-description",
				Price:       price,
			}
			requestBody, err := json.Marshal(&np)
			if err != nil {
				t.Fatal(err)
			}
			req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(requestBody))
			response := tests.ExecuteRequest(req, r)
			tests.AssertStatusCode(t, http.StatusOK, response.Code)
			newNp := Product{
				MerchantId:  merchantId,
				SKUID:       skuId,
				Name:        "new-product-name",
				Description: "new-product-description",
				Price:       100.00,
			}
			newRequestBody, err := json.Marshal(&newNp)
			if err != nil {
				t.Fatal(err)
			}
			newReq, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(newRequestBody))
			newResponse := tests.ExecuteRequest(newReq, r)
			tests.AssertStatusCode(t, http.StatusBadRequest, newResponse.Code)
			message := tests.ParseResponse(t, newResponse)["message"].(string)
			tests.AssertResponseMessage(t, message, "product already exists")
		},
	)

	t.Run(`Given a merchant tries to create a product with a negative price,
         when they make a POST request to the create product endpoint,
         then the API should return a validation error indicating that the price 
         must be a positive number. `,
		func(t *testing.T) {
			merchantId := uuid.New()
			skuId := uuid.New()
			var price float64 = -20.00
			np := Product{
				MerchantId:  merchantId,
				SKUID:       skuId,
				Name:        "some-product-name",
				Description: "some-product-description",
				Price:       price,
			}
			requestBody, err := json.Marshal(&np)
			if err != nil {
				t.Fatal(err)
			}
			req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(requestBody))
			response := tests.ExecuteRequest(req, r)
			tests.AssertStatusCode(t, http.StatusBadRequest, response.Code)
			message := tests.ParseResponse(t, response)["message"].(string)
			tests.AssertResponseMessage(t, message, "price cannot be less than zero")
		},
	)
}

func TestUpdateProduct(t *testing.T) {
	route := "/api/products"
	t.Run("test for invalid json request body",
		func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, route, nil)
			response := tests.ExecuteRequest(req, r)
			tests.AssertStatusCode(t, http.StatusBadRequest, response.Code)
		},
	)
	t.Run(`Given a merchant wants to update a product with valid data,
         When the update product endpoint is called with the correct SKU ID,
         Then the product information should be updated in the database,
         And the response status should be 200 OK. `,
		func(t *testing.T) {
			merchantId := uuid.New()
			skuId := uuid.New()
			var price float64 = 30.00
			np := Product{
				MerchantId:  merchantId,
				SKUID:       skuId,
				Name:        "some-product-name",
				Description: "some-product-description",
				Price:       price,
			}
			createdProductSkuId := createProduct(t, np)
			updatedPrice := 10.01
			updateReq := Product{
				MerchantId:  merchantId,
				SKUID:       createdProductSkuId,
				Name:        "updated-product-name",
				Description: "updated-product-description",
				Price:       updatedPrice,
			}

			requestBody, err := json.Marshal(&updateReq)
			if err != nil {
				t.Fatal(err)
			}
			req, _ := http.NewRequest(http.MethodPatch, route+"/"+createdProductSkuId.String(), bytes.NewBuffer(requestBody))
			response := tests.ExecuteRequest(req, r)
			tests.AssertStatusCode(t, http.StatusOK, response.Code)
			parsedRes := tests.ParseResponse(t, response)
			message := parsedRes["message"].(string)
			tests.AssertResponseMessage(t, message, "product updated successfully")
			data := parsedRes["data"].(map[string]interface{})
			tests.AssertResponseMessage(t, data["merchant_id"].(string), merchantId.String())
			tests.AssertResponseMessage(t, data["sku_id"].(string), createdProductSkuId.String())
			tests.AssertResponseMessage(t, data["name"].(string), updateReq.Name)
			tests.AssertResponseMessage(t, data["description"].(string), updateReq.Description)
			if responsePrice, ok := data["price"].(float64); ok {
				tests.AssertResponseMessage(t, strconv.FormatFloat(responsePrice, 'f', -1, 64), strconv.FormatFloat(updatedPrice, 'f', -1, 64))
			}
		},
	)

	t.Run(`Given a merchant wants to update a product that does not exist,
         When the update product endpoint is called with a non-existent SKU ID,
         Then the endpoint should return a 404 Not Found status,
         And the product information should not be updated in the database. `,
		func(t *testing.T) {
			nonExisitentSkuId := uuid.New()
			updatedPrice := 10.01
			updateReq := Product{
				MerchantId:  uuid.New(),
				SKUID:       nonExisitentSkuId,
				Name:        "updated-product-name",
				Description: "updated-product-description",
				Price:       updatedPrice,
			}

			requestBody, err := json.Marshal(&updateReq)
			if err != nil {
				t.Fatal(err)
			}
			req, _ := http.NewRequest(http.MethodPatch, route+"/"+nonExisitentSkuId.String(), bytes.NewBuffer(requestBody))
			response := tests.ExecuteRequest(req, r)
			tests.AssertStatusCode(t, http.StatusNotFound, response.Code)
			parsedRes := tests.ParseResponse(t, response)
			message := parsedRes["message"].(string)
			tests.AssertResponseMessage(t, message, "product not found")
		},
	)
}

func createProduct(t testing.TB, np Product) uuid.UUID {
	t.Helper()
	requestBody, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(requestBody))
	response := tests.ExecuteRequest(req, r)
	data := tests.ParseResponse(t, response)["data"].(map[string]interface{})
	if id, ok := data["sku_id"].(string); !ok && id == np.SKUID.String() {
		t.Fatal()
	}

	return np.SKUID
}

type Product struct {
	SKUID       uuid.UUID `json:"sku_id"`
	MerchantId  uuid.UUID `json:"merchant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
}
