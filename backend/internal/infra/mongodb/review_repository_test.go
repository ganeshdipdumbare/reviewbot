package mongodb

import (
	"backend/internal/review"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ReviewTestSuite struct {
	suite.Suite
	mongoContainer *mongoContainer
}

func (suite *ReviewTestSuite) SetupSuite() {
	ctx := context.Background()
	mongoContainer, err := getMongoTestContainer(ctx)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.mongoContainer = mongoContainer
}

func (suite *ReviewTestSuite) TearDownSuite() {
	suite.mongoContainer.Container.Terminate(context.Background())
}

func TestReviewTestSuite(t *testing.T) {
	suite.Run(t, new(ReviewTestSuite))
}

func (suite *ReviewTestSuite) TestNewClient() {
	client, err := NewClient(suite.mongoContainer.URI, "test")
	suite.NoError(err)
	suite.NotNil(client)
}

func (suite *ReviewTestSuite) TestNewReviewRepository() {
	client, err := NewClient(suite.mongoContainer.URI, "test")
	suite.NoError(err)
	suite.NotNil(client)

	repo, err := NewReviewRepository(client, "test")
	suite.NoError(err)
	suite.NotNil(repo)
}

func (suite *ReviewTestSuite) TestUpsert() {
	client, err := NewClient(suite.mongoContainer.URI, "test")
	suite.NoError(err)
	suite.NotNil(client)

	repo, err := NewReviewRepository(client, "test")
	suite.NoError(err)
	suite.NotNil(repo)

	r := &review.Review{
		UserID:    "userid",
		ProductID: "productid",
		Body:      "body",
		Rating:    5,
	}

	resp, err := repo.Upsert(r)
	suite.NoError(err)
	suite.NotNil(resp)
}

func (suite *ReviewTestSuite) TestGet() {
	client, err := NewClient(suite.mongoContainer.URI, "test")
	suite.NoError(err)
	suite.NotNil(client)

	repo, err := NewReviewRepository(client, "test")
	suite.NoError(err)
	suite.NotNil(repo)

	r := &review.Review{
		UserID:    "userid",
		ProductID: "productid",
		Body:      "body",
		Rating:    5,
	}

	resp, err := repo.Upsert(r)
	suite.NoError(err)
	suite.NotNil(resp)

	review, err := repo.Get(resp.ID)
	suite.NoError(err)
	suite.NotNil(review)
}
