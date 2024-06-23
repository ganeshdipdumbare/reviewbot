package simplenlp_test

import (
	"testing"

	"backend/internal/infra/msgnlp"
	"backend/internal/infra/msgnlp/simplenlp"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestPredictIntent(t *testing.T) {
	nlp := simplenlp.NewSimpleNLP()

	tests := []struct {
		name     string
		request  *msgnlp.MessageIntentRequest
		expected *msgnlp.MessageIntentResponse
	}{
		{
			name: "Test Greetings Intent",
			request: &msgnlp.MessageIntentRequest{
				Text: "Hello, how are you?",
			},
			expected: &msgnlp.MessageIntentResponse{
				Intent: msgnlp.IntentGreet,
			},
		},
		{
			name: "Test Submit Review Intent",
			request: &msgnlp.MessageIntentRequest{
				Text: "Please provide your review.",
			},
			expected: &msgnlp.MessageIntentResponse{
				Intent: msgnlp.IntentSubmitReview,
			},
		},
		{
			name: "Test Submit Rating Intent",
			request: &msgnlp.MessageIntentRequest{
				Text: "I rate 5 stars.",
			},
			expected: &msgnlp.MessageIntentResponse{
				Intent: msgnlp.IntentSubmitRating,
			},
		},
		{
			name: "Test Get Product Intent",
			request: &msgnlp.MessageIntentRequest{
				Text: "Tell me more about the product iPhone.",
			},
			expected: &msgnlp.MessageIntentResponse{
				Intent: msgnlp.IntentGetProduct,
				Entities: map[string]string{
					"product": "iphone",
				},
			},
		},
		{
			name: "Test Goodbye Intent",
			request: &msgnlp.MessageIntentRequest{
				Text: "Goodbye, see you later!",
			},
			expected: &msgnlp.MessageIntentResponse{
				Intent: msgnlp.IntentGoodbye,
			},
		},
		{
			name: "Test System Start Review Intent",
			request: &msgnlp.MessageIntentRequest{
				Text: "system@startreview",
			},
			expected: &msgnlp.MessageIntentResponse{
				Intent: msgnlp.IntentSystemStartReview,
			},
		},
		{
			name: "Test Unknown Intent",
			request: &msgnlp.MessageIntentRequest{
				Text: "This is an unknown intent.",
			},
			expected: &msgnlp.MessageIntentResponse{
				Intent: msgnlp.IntentUnknown,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := nlp.PredictIntent(test.request)
			assert.NoError(t, err)
			if !cmp.Equal(response, test.expected) {
				t.Errorf("Response not as expected: %s", cmp.Diff(response, test.expected))
			}
		})
	}
}
