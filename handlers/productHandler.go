package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/adityagoel/product-api/data"
	"github.com/gorilla/mux"
)

type Products struct {
	thisLogger *log.Logger
}

func NewProducts(thisLogger *log.Logger) *Products {
	return &Products{thisLogger}
}

func (h *Products) UpdateSingleProduct(responseWriter http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)

	id, errorWhileTypeCasting := strconv.Atoi(vars["id"])

	if errorWhileTypeCasting != nil {
		http.Error(responseWriter, "Unable to typeCast the request-parameter", http.StatusBadRequest)
		return
	}

	h.thisLogger.Println("Handle PUT Products")

	thisProd := request.Context().Value(KeyProduct{}).(data.Product)

	// Adding this prooduct to to our Temporary DataStore Now.
	err := data.UpdateSingleProduct(id, &thisProd)

	if err != nil {
		http.Error(responseWriter, "Unable to update this product to DataStore", http.StatusBadRequest)
	}

	h.thisLogger.Printf("Prod: %#v", thisProd)
}

func (h *Products) GetProducts(responseWriter http.ResponseWriter, request *http.Request) {

	h.thisLogger.Println("Handle Get Products")

	listOfProducts := data.GetProducts()
	errorWhileEncoding := listOfProducts.ToJson(responseWriter)

	if errorWhileEncoding != nil {
		http.Error(responseWriter, "Unable to encode the JSON", http.StatusInternalServerError)
	}
}

func (h *Products) AddProduct(responseWriter http.ResponseWriter, request *http.Request) {

	h.thisLogger.Println("Handle Post Product")

	thisProd := request.Context().Value(KeyProduct{}).(data.Product)

	// Adding this prooduct to tour Temporary DataStore Now.
	data.AddSingleProduct(&thisProd)

	h.thisLogger.Printf("Prod: %#v", thisProd)
}

type KeyProduct struct{}

func (h *Products) ValidateProducts(thisProductsHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		thisProd := data.Product{}

		errorWhileUnMarshalling := thisProd.FromJson(request.Body)

		if errorWhileUnMarshalling != nil {
			http.Error(responseWriter, "Unable to unmarshall the JSON", http.StatusBadRequest)
			h.thisLogger.Println("Error while unmarshalling and reading the product from JSON.")
			return
		}
		h.thisLogger.Println("JSON has been succesfully unmarshalled and product has been read from it.")

		thisContext := context.WithValue(request.Context(), KeyProduct{}, thisProd)
		request = request.WithContext(thisContext)

		thisProductsHandler.ServeHTTP(responseWriter, request)

	})
}
