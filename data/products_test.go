package data

import "testing"

func TestValidation(t *testing.T) {
	p := &Product{
		Name: "Hello",
		Price: 0.99,
		SKU: "abc-abc-s",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
