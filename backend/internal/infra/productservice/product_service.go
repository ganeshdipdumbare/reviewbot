package productservice

type Product struct {
	ID          string
	Name        string
	Description string
	Price       string
	Rating      int
}

// ProductService provides access to the product storage.
//
//go:generate mockgen -destination ../../mocks/productservice/mock_product_service.go -package=productservice backend/internal/infra/productservice ProductService
type ProductService interface {
	GetProduct(id string) (*Product, error)
}
