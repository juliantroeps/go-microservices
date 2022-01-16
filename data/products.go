package data

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

// Product defines the structure for an API products
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"desc"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// Products represents a list of products
type Products []*Product

// ToJSON converts the list of products to encoded json
func (p *Products) ToJSON(w io.Writer) error {
	// Create a new Encode
	encoder := json.NewEncoder(w)

	// Return encoded JSON
	return encoder.Encode(p)
}

// FromJSON converts JSON input to a Product instance
func (p *Product) FromJSON(r io.Reader) error {
	// Create a new Decoder
	decoder := json.NewDecoder(r)

	// Return decoded List
	return decoder.Decode(p)
}

// Our "database"
var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Milk coffee with milk and coffee.",
		Price:       2.45,
		SKU:         "PA0001",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk.",
		Price:       1.99,
		SKU:         "PA0002",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

// GetProducts gets all products
func GetProducts() Products {
	return productList
}

// AddProduct adds a new product
func AddProduct(p *Product) {
	p.ID = getNextID()
	p.SKU = "PA00" + strconv.Itoa(getNextID())
	p.CreatedOn = time.Now().UTC().String()
	p.UpdatedOn = time.Now().UTC().String()

	productList = append(productList, p)
}

// UpdateProduct updates a product by its id
func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	p.UpdatedOn = time.Now().UTC().String()
	productList[pos] = p

	return nil
}

// === === === HELPER  === === === ///

var ErrProductNotFound = fmt.Errorf("product not found")

// findProduct Get the product by id from the data store
func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
}

// getNextID simply bumps the id
func getNextID() int {
	lastItem := productList[len(productList)-1]
	return lastItem.ID + 1
}
