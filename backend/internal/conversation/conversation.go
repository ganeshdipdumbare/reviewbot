package conversation

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidField = errors.New("invalid field")
)

type Conversation struct {
	ID        string
	CreatedAt time.Time
	UserID    string
	Messages  []Message
}

type Message struct {
	Text      string
	Timestamp time.Time
}

// Validate validates the conversation and returns an error if it is not valid.
func (c *Conversation) Validate() error {
	if c.UserID == "" {
		return fmt.Errorf("%w: user_id", ErrInvalidField)
	}
	if len(c.Messages) == 0 {
		return fmt.Errorf("%w: messages", ErrInvalidField)
	}
	return nil
}

func (c *Conversation) AddMessage(text string) {
	c.Messages = append(c.Messages, Message{
		Text:      text,
		Timestamp: time.Now(),
	})
}

// Repository provides access to the conversation storage.
//
//go:generate mockgen -destination ../mocks/convrepo/mock_conversation_repository.go -package=convrepo backend/internal/conversation Repository
type Repository interface {
	Get(id string) (*Conversation, error)
	Upsert(c *Conversation) (*Conversation, error)
}
