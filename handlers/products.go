package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"go-microservices/data"
	"log"
	"net/http"
	"strconv"
)

// Products is a http.Handler and defines the handler structure
type Products struct {
	logger *log.Logger
}

// NewProducts initializes a products' handler instance with given logger
func NewProducts(logger *log.Logger) *Products {
	return &Products{
		logger: logger,
	}
}

// GetProducts returns all products from the data store
func (products *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	products.logger.Println("Handle GET products")

	// Fetch the products from the data store
	listOfProducts := data.GetProducts()

	// Set content type
	rw.Header().Set("Content-Type", "application/json")

	// Serialize the list to JSON
	err := listOfProducts.ToJSON(rw)

	// Handle error
	if err != nil {
		http.Error(rw, "Unable to encode json", http.StatusInternalServerError)
	}
}

// AddProduct adds a new item
func (products *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	products.logger.Println("Handle POST products")

	// Get the product data from the request context
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	// Add product to data store
	data.AddProduct(prod)
}

// UpdateProduct updates a given item
func (products *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	// Get variables (from path)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		products.logger.Print(err)
		http.Error(rw, "Error while formatting id", http.StatusBadRequest)
	}

	// Log
	products.logger.Println("Handle PUT products", id)

	// Get the product data from the request context
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found.", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found.", http.StatusInternalServerError)
		return
	}
}

// KeyProduct Create a placeholder
type KeyProduct struct{}

// MiddlewareProductValidation is a custom middleware for validation products
func (products Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// New data.Product instance
		prod := &data.Product{}

		// Serialize the JSON to product
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal to JSON", http.StatusBadRequest)
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
