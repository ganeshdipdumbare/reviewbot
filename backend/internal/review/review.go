package review

import (
	"errors"
	"fmt"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusClosed  Status = "closed"
)

var (
	ErrInvalidField = errors.New("invalid field")
)

// Review represents a review of a product
type Review struct {
	ID             string
	CreatedAt      int64
	ProductID      string
	UserID         string
	Body           string
	Rating         int
	ConversationID string
	Status         Status
}

func (r *Review) Validate() error {
	if r.ProductID == "" {
		return fmt.Errorf("%w: product_id", ErrInvalidField)
	}
	if r.UserID == "" {
		return fmt.Errorf("%w: user_id", ErrInvalidField)
	}
	if r.Body == "" {
		return fmt.Errorf("%w: body", ErrInvalidField)
	}
	if r.Rating < 1 || r.Rating > 5 {
		return fmt.Errorf("%w: rating", ErrInvalidField)
	}
	if r.ConversationID == "" {
		return fmt.Errorf("%w: conversation_id", ErrInvalidField)
	}
	if r.Status == "" {
		return fmt.Errorf("%w: status", ErrInvalidField)
	}
	return nil
}

// Repository provides access to the review storage.
//
//go:generate mockgen -destination ../mocks/reviewrepo/mock_review_repository.go -package=reviewrepo backend/internal/review Repository
type Repository interface {
	Get(id string) (*Review, error)
	Upsert(r *Review) (*Review, error)
}
