package app

import (
	"backend/internal/conversation"
	"backend/internal/infra/msgnlp"
	"backend/internal/review"
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

func (a *conversationApp) Converse(req *ConverseRequest) (*ConverseResponse, error) {
	// check if the conversation exists, if not create a new conversation
	convo, err := a.getOrCreateConversation(req)
	if err != nil {
		a.logger.Debug("failed to get or create conversation", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get or create conversation: %w", err)
	}

	conversationResponse := &ConverseResponse{
		ConversationID: req.ConversationID,
		ProductID:      req.ProductID,
		UserID:         req.UserID,
		ReviewID:       req.ReviewID,
		Text:           "",
	}

	// Get message intent
	intent, err := a.msgNLP.PredictIntent(&msgnlp.MessageIntentRequest{
		Text: req.Text,
	})
	if err != nil {
		a.logger.Debug("failed to get message intent", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get message intent: %w", err)
	}

	// Handle different intents and create appropriate response message
	switch intent.Intent {
	case msgnlp.IntentSystemInitiateConversation:
		// create or update conversation
		a.logger.Debug("initiating conversation")
		conversationResponse.ConversationID = convo.ID
		conversationResponse.Text = "Hello! How can I help you today?"
	case msgnlp.IntentGoodbye:
		// end conversation
		a.logger.Debug("ending conversation")
		err := a.EndConversation(&EndConversationRequest{
			ConversationID: req.ConversationID,
			ReviewID:       req.ReviewID,
		})
		if err != nil {
			a.logger.Debug("failed to end conversation", slog.Any("error", err))
			return nil, fmt.Errorf("failed to end conversation: %w", err)
		}
		conversationResponse.Text = "Goodbye! Dont hesitate to reach out if you need help."
	case msgnlp.IntentGreet:
		a.logger.Debug("greeting")
		conversationResponse.Text = "Hello!"
	case msgnlp.IntentSystemStartReview:
		a.logger.Debug("starting review")
		// get conversation and update text
		conversationResponse.ConversationID = convo.ID
		conversationResponse.Text = "Please provide your review for iPhone 13."
		// create review
		rev := &review.Review{
			CreatedAt:      time.Now().Unix(),
			ProductID:      req.ProductID,
			UserID:         req.UserID,
			Body:           req.Text,
			Rating:         0,
			ConversationID: convo.ID,
			Status:         review.StatusPending,
		}
		rev, err = a.reviewRepo.Upsert(rev)
		if err != nil {
			a.logger.Debug("failed to create review", slog.Any("error", err))
			return nil, fmt.Errorf("failed to create review: %w", err)
		}
		conversationResponse.ReviewID = rev.ID
	case msgnlp.IntentSubmitReview:
		a.logger.Debug("submitting review")
		// update review
		review, err := a.reviewRepo.Get(req.ReviewID)
		if err != nil {
			return nil, fmt.Errorf("failed to get review: %w", err)
		}
		review.Body = req.Text
		_, err = a.reviewRepo.Upsert(review)
		if err != nil {
			a.logger.Debug("failed to update review", slog.Any("error", err))
			return nil, fmt.Errorf("failed to update review: %w", err)
		}
		conversationResponse.Text = "Thank you for your review! Please provide a rating from 1 to 5."
	case msgnlp.IntentSubmitRating:
		a.logger.Debug("submitting rating")
		// update review
		rev, err := a.reviewRepo.Get(req.ReviewID)
		if err != nil {
			a.logger.Debug("failed to get review", slog.Any("error", err))
			return nil, fmt.Errorf("failed to get review: %w", err)
		}
		rev.Rating, err = strconv.Atoi(req.Text)
		if err != nil {
			a.logger.Debug("failed to convert rating to integer", slog.Any("error", err))
			return nil, fmt.Errorf("failed to convert rating to integer: %w", err)
		}
		rev.Status = review.StatusClosed
		_, err = a.reviewRepo.Upsert(rev)
		if err != nil {
			a.logger.Debug("failed to update review", slog.Any("error", err))
			return nil, fmt.Errorf("failed to update review: %w", err)
		}
		conversationResponse.Text = "Thank you for your rating! Your review has been submitted."
	case msgnlp.IntentGetProduct:
		a.logger.Debug("getting product")
		// get product
		product, err := a.productService.GetProduct(req.Text)
		if err != nil {
			a.logger.Debug("failed to get product", slog.Any("error", err))
			return nil, fmt.Errorf("failed to get product: %w", err)
		}
		conversationResponse.Text = fmt.Sprintf("Product: %s\nDescription: %s", product.Name, product.Description)
	default:
		a.logger.Debug("unknown intent")
		conversationResponse.Text = "I'm sorry, I don't understand."
	}

	a.logger.Debug("responding", slog.Any("response", conversationResponse.Text))
	convo.AddMessage(conversationResponse.Text)
	a.logger.Debug("updating conversation", slog.Any("conversation", *convo))
	_, err = a.convoRepo.Upsert(convo)
	if err != nil {
		a.logger.Debug("failed to update conversation", slog.Any("error", err))
		return nil, fmt.Errorf("failed to update conversation: %w", err)
	}
	return conversationResponse, nil
}

// getOrCreateConversation gets or creates a conversation
// based on the request conversation ID.
func (a *conversationApp) getOrCreateConversation(req *ConverseRequest) (*conversation.Conversation, error) {
	var (
		convo *conversation.Conversation
		err   error
	)

	if req.ConversationID == "" {
		convo, err = a.convoRepo.Upsert(&conversation.Conversation{
			UserID:    req.UserID,
			CreatedAt: time.Now(),
			Messages: []conversation.Message{
				{
					Text:      req.Text,
					Timestamp: time.Now(),
				},
			},
		})
		if err != nil {
			a.logger.Debug("failed to create conversation", slog.Any("error", err))
			return nil, fmt.Errorf("failed to create conversation: %w", err)
		}
	} else {
		convo, err = a.convoRepo.Get(req.ConversationID)
		if err != nil {
			a.logger.Debug("failed to get conversation", slog.Any("error", err))
			return nil, fmt.Errorf("failed to get conversation: %w", err)
		}
		convo.AddMessage(req.Text)
	}
	return convo, nil
}
