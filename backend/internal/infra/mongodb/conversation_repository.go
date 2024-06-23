package mongodb

import (
	"backend/internal/conversation"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type conversationRepository struct {
	client         *Client
	collectionName string
	collection     *mongo.Collection
}

// NewConversationRepository creates a new conversation repository
func NewConversationRepository(client *Client, collectionName string) (*conversationRepository, error) {
	if client == nil {
		return nil, fmt.Errorf("%w: client", ErrEmptyArg)
	}
	if collectionName == "" {
		return nil, fmt.Errorf("%w: collectionName", ErrEmptyArg)
	}
	conversationCollection := client.mdbclient.Database(client.dbName).Collection(collectionName)
	return &conversationRepository{
		client:         client,
		collectionName: collectionName,
		collection:     conversationCollection,
	}, nil
}

type Conversation struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UserID    string             `bson:"user_id"`
	Messages  []Message          `bson:"messages"`
}

type Message struct {
	Text      string    `bson:"text"`
	Timestamp time.Time `bson:"timestamp"`
}

func convertToMongoConversation(c *conversation.Conversation) (*Conversation, error) {
	var messages []Message
	for _, m := range c.Messages {
		messages = append(messages, Message{
			Text:      m.Text,
			Timestamp: m.Timestamp,
		})
	}

	conversation := &Conversation{
		CreatedAt: c.CreatedAt,
		UserID:    c.UserID,
		Messages:  messages,
	}

	if c.ID != "" {
		id, err := primitive.ObjectIDFromHex(c.ID)
		if err != nil {
			return nil, err
		}
		conversation.ID = id
	} else {
		conversation.ID = primitive.NewObjectID()
	}
	return conversation, nil
}

// Get retrieves a conversation by its ID
func (r *conversationRepository) Get(id string) (*conversation.Conversation, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", objectID}}
	var result Conversation
	err = r.collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	var messages []conversation.Message
	for _, m := range result.Messages {
		messages = append(messages, conversation.Message{
			Text:      m.Text,
			Timestamp: m.Timestamp,
		})
	}

	return &conversation.Conversation{
		ID:        result.ID.Hex(),
		CreatedAt: result.CreatedAt,
		UserID:    result.UserID,
		Messages:  messages,
	}, nil
}

// Upsert creates a new conversation if the conversation ID is empty,
func (r *conversationRepository) Upsert(c *conversation.Conversation) (*conversation.Conversation, error) {
	conversation, err := convertToMongoConversation(c)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", conversation.ID}}
	update := bson.D{{"$set", conversation}}
	_, err = r.collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	c.ID = conversation.ID.Hex()
	return c, nil
}
