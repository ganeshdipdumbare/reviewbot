package mongodb

import (
	"context"
	"fmt"

	"backend/internal/review"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type reviewRepository struct {
	client         *Client
	collectionName string
	collection     *mongo.Collection
}

// NewReviewRepository creates a new review repository
func NewReviewRepository(client *Client, collectionName string) (*reviewRepository, error) {
	if client == nil {
		return nil, fmt.Errorf("%w: client", ErrEmptyArg)
	}
	if collectionName == "" {
		return nil, fmt.Errorf("%w: collectionName", ErrEmptyArg)
	}
	reviewCollection := client.mdbclient.Database(client.dbName).Collection(collectionName)
	return &reviewRepository{
		client:         client,
		collectionName: collectionName,
		collection:     reviewCollection,
	}, nil
}

type Review struct {
	ID             primitive.ObjectID `bson:"_id"`
	CreatedAt      int64              `bson:"created_at"`
	ProductID      string             `bson:"product_id"`
	UserID         string             `bson:"user_id"`
	Body           string             `bson:"body"`
	Rating         int                `bson:"rating"`
	ConversationID string             `bson:"conversation_id"`
	Status         string             `bson:"status"`
}

func convertToMongoReview(r *review.Review) (*Review, error) {
	review := &Review{
		CreatedAt:      r.CreatedAt,
		ProductID:      r.ProductID,
		UserID:         r.UserID,
		Body:           r.Body,
		Rating:         r.Rating,
		ConversationID: r.ConversationID,
		Status:         string(r.Status),
	}

	if r.ID != "" {
		id, err := primitive.ObjectIDFromHex(r.ID)
		if err != nil {
			return nil, err
		}
		review.ID = id
	} else {
		review.ID = primitive.NewObjectID()
	}
	return review, nil
}

// Get retrieves a review by its ID
func (r *reviewRepository) Get(id string) (*review.Review, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", objectID}}
	var result Review
	err = r.collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	review := &review.Review{
		ID:             result.ID.Hex(),
		CreatedAt:      result.CreatedAt,
		ProductID:      result.ProductID,
		UserID:         result.UserID,
		Body:           result.Body,
		Rating:         result.Rating,
		ConversationID: result.ConversationID,
		Status:         review.Status(result.Status),
	}
	return review, nil
}

// Upsert creates a new review if the review ID is empty,
func (r *reviewRepository) Upsert(c *review.Review) (*review.Review, error) {
	review, err := convertToMongoReview(c)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", review.ID}}
	update := bson.D{
		{"$set", review},
	}
	_, err = r.collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}
	c.ID = review.ID.Hex()
	return c, nil
}
