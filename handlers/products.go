package handlers

import (
	"go-microservices/data"
	"log"
	"net/http"
	"regexp"
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

// ServeHTTP is the main entry point for the handler and satisfies the http.Handler (overwrite the standard method)
// ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL.
func (products *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// Get all
	if r.Method == http.MethodGet {
		products.getProducts(rw, r)
		return
	}

	// Add item
	if r.Method == http.MethodPost {
		products.addProduct(rw, r)
		return
	}

	// Update item
	if r.Method == http.MethodPut {
		// Get the url path
		path := r.URL.Path

		// Log
		products.logger.Println("PUT", path)

		// Find with regexp
		regex := regexp.MustCompile(`\/([0-9]+)`)
		group := regex.FindAllStringSubmatch(path, -1)

		if len(group) != 1 {
			products.logger.Println("More than one id in path", path)
			http.Error(rw, "Invalid uri (1)", http.StatusBadRequest)
			return
		}

		if len(group[0]) != 2 {
			products.logger.Println("More than one capture group in path", path)
			http.Error(rw, "Invalid uri (2)", http.StatusBadRequest)
			return
		}

		// Get the id as strong
		idString := group[0][1]

		// Convert to int
		id, err := strconv.Atoi(idString)

		// Handle err
		if err != nil {
			products.logger.Println("Unable to convert to number", idString)
			http.Error(rw, "Invalid uri (3)", http.StatusBadRequest)
			return
		}

		products.updateProduct(id, rw, r)

		return
	}

	// Catch all
	// If no method is satisfied return an error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// getProducts returns all products from the data store
func (products *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
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

// addProduct adds a new item
func (products *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	products.logger.Println("Handle POST products")

	// New data.Product instance
	newProduct := &data.Product{}

	// Serialize the JSON to product
	err := newProduct.FromJSON(r.Body)
	if err != nil {
		products.logger.Print(err)
		http.Error(rw, "Unable to unmarshal to JSON", http.StatusBadRequest)
	}

	data.AddProduct(newProduct)
}

// updateProduct updates a given item
func (products *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	products.logger.Println("Handle PUT products")

	// New data.Product instance
	newProduct := &data.Product{}

	// Serialize the JSON to product
	err := newProduct.FromJSON(r.Body)
	if err != nil {
		products.logger.Print(err)
		http.Error(rw, "Unable to unmarshal to JSON", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, newProduct)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found.", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found.", http.StatusInternalServerError)
		return
	}
}
