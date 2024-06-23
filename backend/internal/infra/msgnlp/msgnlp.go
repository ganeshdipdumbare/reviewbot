package msgnlp

type Intent string

const (
	IntentUnknown                    Intent = "unknown"
	IntentGreet                      Intent = "greet"
	IntentSubmitReview               Intent = "submit_review"
	IntentSubmitRating               Intent = "submit_rating"
	IntentGetProduct                 Intent = "get_product"
	IntentGoodbye                    Intent = "goodbye"
	IntentSystemStartReview          Intent = "system_start_review"
	IntentSystemInitiateConversation Intent = "system_initiate_conversation"
)

type MessageIntentRequest struct {
	Text string
}

type MessageIntentResponse struct {
	Intent   Intent
	Entities map[string]string
}

// MessageIntentService is the interface for the message intent service
//
//go:generate mockgen -destination ../../mocks/msgnlp/mock_msgnlp.go -package=msgnlp . MessageIntentService
type MessageIntentService interface {
	PredictIntent(req *MessageIntentRequest) (*MessageIntentResponse, error)
}
