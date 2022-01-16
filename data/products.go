package data

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"regexp"
	"strconv"
	"time"
)

// Product defines the structure for an API products
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"desc"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// Validate the product fields
func (p *Product) Validate() error {
	// Create new validator
	validate := validator.New()

	// Custom validation for SKU
	err := validate.RegisterValidation("sku", func(fl validator.FieldLevel) bool {
		// Format: abc-abc-abc
		re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
		matches := re.FindAllString(fl.Field().String(), -1)

		if len(matches) != 1 {
			return false
		}

		return true
	})

	if err != nil {
		return err
	}

	return validate.Struct(p)
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
