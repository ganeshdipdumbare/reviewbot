package mongodb

import (
	"backend/internal/conversation"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ConversationTestSuite struct {
	suite.Suite
	mongoContainer *mongoContainer
}

func (suite *ConversationTestSuite) SetupSuite() {
	ctx := context.Background()
	mongoContainer, err := getMongoTestContainer(ctx)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.mongoContainer = mongoContainer
}

func (suite *ConversationTestSuite) TearDownSuite() {
	suite.mongoContainer.Container.Terminate(context.Background())
}

func TestConversationTestSuite(t *testing.T) {
	suite.Run(t, new(ConversationTestSuite))
}

func (suite *ConversationTestSuite) TestNewClient() {
	client, err := NewClient(suite.mongoContainer.URI, "test")
	suite.NoError(err)
	suite.NotNil(client)
}

func (suite *ConversationTestSuite) TestNewConversationRepository() {
	client, err := NewClient(suite.mongoContainer.URI, "test")
	suite.NoError(err)
	suite.NotNil(client)

	repo, err := NewConversationRepository(client, "test")
	suite.NoError(err)
	suite.NotNil(repo)
}

func (suite *ConversationTestSuite) TestUpsert() {
	client, err := NewClient(suite.mongoContainer.URI, "test")
	suite.NoError(err)
	suite.NotNil(client)

	repo, err := NewConversationRepository(client, "test")
	suite.NoError(err)
	suite.NotNil(repo)

	c := &conversation.Conversation{
		UserID:    "userid",
		CreatedAt: time.Now(),
		Messages: []conversation.Message{
			{
				Text:      "hello",
				Timestamp: time.Now(),
			},
		},
	}
	resp, err := repo.Upsert(c)
	suite.NoError(err)
	suite.NotNil(resp)
}

func (suite *ConversationTestSuite) TestGet() {
	client, err := NewClient(suite.mongoContainer.URI, "test")
	suite.NoError(err)
	suite.NotNil(client)

	repo, err := NewConversationRepository(client, "test")
	suite.NoError(err)
	suite.NotNil(repo)

	c := &conversation.Conversation{
		UserID:    "userid",
		CreatedAt: time.Now(),
		Messages: []conversation.Message{
			{
				Text:      "hello",
				Timestamp: time.Now(),
			},
		},
	}
	resp, err := repo.Upsert(c)
	suite.NoError(err)
	suite.Assert().NotEmpty(resp.ID)

	convo, err := repo.Get(resp.ID)
	suite.NoError(err)
	suite.Equal(resp.ID, convo.ID)
	suite.Equal(resp.UserID, convo.UserID)
	suite.Equal(resp.CreatedAt.Unix(), convo.CreatedAt.Unix())
}
