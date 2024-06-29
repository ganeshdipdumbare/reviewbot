package simplenlp

import (
	"strings"

	"backend/internal/infra/msgnlp"
)

type simpleNLP struct{}

// NewSimpleNLP creates a new simple NLP service
func NewSimpleNLP() *simpleNLP {
	return &simpleNLP{}
}

// GetIntent function to determine the intent of the input
func (s *simpleNLP) PredictIntent(
	req *msgnlp.MessageIntentRequest,
) (*msgnlp.MessageIntentResponse, error) {
	// convert input text to lowercase and remove punctuation
	input := strings.ToLower(req.Text)
	input = strings.ReplaceAll(input, ",", "")
	input = strings.ReplaceAll(input, ".", "")
	input = strings.ReplaceAll(input, "?", "")
	input = strings.ReplaceAll(input, "!", "")

	words := strings.Fields(input)

	// Define keywords for each intent
	greetings := []string{"hello", "hi", "hey", "greetings"}
	submitReview := []string{"review", "feedback", "comment"}
	submitRating := []string{"1", "2", "3", "4", "5"}
	getProduct := []string{"product", "item", "details", "info"}
	byeBye := []string{"bye", "goodbye", "see you", "farewell"}
	systemInitiateConversation := []string{"system@initiateconversation"}
	systemStartReview := []string{"system@startreview"}

	// Check for each intent
	resp := &msgnlp.MessageIntentResponse{
		Intent: msgnlp.IntentUnknown,
	}
	for _, word := range words {
		if containsIgnoreCase(systemInitiateConversation, word) {
			resp.Intent = msgnlp.IntentSystemInitiateConversation
		} else if containsIgnoreCase(greetings, word) {
			resp.Intent = msgnlp.IntentGreet
		} else if containsIgnoreCase(submitReview, word) {
			resp.Intent = msgnlp.IntentSubmitReview
		} else if containsIgnoreCase(submitRating, word) {
			resp.Intent = msgnlp.IntentSubmitRating
		} else if containsIgnoreCase(getProduct, word) {
			resp.Intent = msgnlp.IntentGetProduct
			// add last word as product entity
			resp.Entities = map[string]string{
				"product": words[len(words)-1],
			}
		} else if containsIgnoreCase(byeBye, word) {
			resp.Intent = msgnlp.IntentGoodbye
		} else if containsIgnoreCase(systemStartReview, word) {
			resp.Intent = msgnlp.IntentSystemStartReview
		}
	}

	return resp, nil
}

// containsIgnoreCase helper function to check if a slice contains a given word (case-insensitive)
func containsIgnoreCase(slice []string, word string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, word) {
			return true
		}
	}
	return false
}
