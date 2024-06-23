package simpleproduct

import (
	"backend/internal/infra/productservice"
)

type simpleProductService struct{}

// NewSimpleProductService creates a new simple product service
func NewSimpleProductService() *simpleProductService {
	return &simpleProductService{}
}

// GetProduct retrieves a product by its ID
func (s *simpleProductService) GetProduct(id string) (*productservice.Product, error) {
	return &productservice.Product{
		ID:          id,
		Name:        "Iphone 13",
		Description: "The latest iPhone model, 16GB RAM, 256GB storage",
		Price:       "999.99 $",
		Rating:      5,
	}, nil
}
