package simpleproduct_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/internal/infra/productservice"
	"backend/internal/infra/productservice/simpleproduct"
)

func TestSimpleProductService_GetProduct(t *testing.T) {
	service := simpleproduct.NewSimpleProductService()

	// Test case 1: Valid product ID
	productID := "123"
	expectedProduct := &productservice.Product{
		ID:          productID,
		Name:        "Iphone 13",
		Description: "The latest iPhone model, 16GB RAM, 256GB storage",
		Price:       "999.99 $",
		Rating:      5,
	}

	product, err := service.GetProduct(productID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
}
