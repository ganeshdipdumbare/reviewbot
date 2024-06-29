package app

import (
	"errors"
	"fmt"
	"log/slog"

	"backend/internal/conversation"
	"backend/internal/infra/msgnlp"
	"backend/internal/infra/productservice"
	"backend/internal/review"
)

var ErrEmptyArg = errors.New("empty argument")

type ConverseRequest struct {
	ConversationID string
	ReviewID       string
	UserID         string
	ProductID      string
	Text           string
}

type ConverseResponse struct {
	ConversationID string
	ReviewID       string
	UserID         string
	ProductID      string
	Text           string
}

type EndConversationRequest struct {
	ConversationID string
	ReviewID       string
}

// ConversationApp is the interface for the conversation app
//
//go:generate mockgen -destination ../mocks/app/mock_conversation_app.go -package=app backend/internal/app ConversationApp
type ConversationApp interface {
	Converse(req *ConverseRequest) (*ConverseResponse, error)
	EndConversation(req *EndConversationRequest) error
}

type conversationApp struct {
	logger         *slog.Logger
	reviewRepo     review.Repository
	convoRepo      conversation.Repository
	msgNLP         msgnlp.MessageIntentService
	productService productservice.ProductService
}

func NewConversationApp(logger *slog.Logger, reviewRepo review.Repository, convoRepo conversation.Repository, msgNLP msgnlp.MessageIntentService, productService productservice.ProductService) (*conversationApp, error) {
	if reviewRepo == nil {
		return nil, fmt.Errorf("%w: reviewRepo", ErrEmptyArg)
	}
	if convoRepo == nil {
		return nil, fmt.Errorf("%w: convoRepo", ErrEmptyArg)
	}
	if msgNLP == nil {
		return nil, fmt.Errorf("%w: msgNLP", ErrEmptyArg)
	}
	if productService == nil {
		return nil, fmt.Errorf("%w: productService", ErrEmptyArg)
	}

	return &conversationApp{
		logger:         logger,
		reviewRepo:     reviewRepo,
		convoRepo:      convoRepo,
		msgNLP:         msgNLP,
		productService: productService,
	}, nil
}
