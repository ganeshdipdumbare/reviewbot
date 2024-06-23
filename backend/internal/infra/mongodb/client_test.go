package mongodb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MongoTestSuite struct {
	mongoContainer *mongoContainer
	suite.Suite
}

func (suite *MongoTestSuite) SetupSuite() {
	ctx := context.Background()
	mongoContainer, err := getMongoTestContainer(ctx)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.mongoContainer = mongoContainer
}

func (suite *MongoTestSuite) TearDownSuite() {
	suite.mongoContainer.Container.Terminate(context.Background())
}

func TestMongoTestSuite(t *testing.T) {
	suite.Run(t, new(MongoTestSuite))
}

type mongoContainer struct {
	Container testcontainers.Container
	Ip        string
	Port      string
	URI       string
}

func getMongoTestContainer(ctx context.Context) (*mongoContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:7.0.3",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForHTTP("/").WithStartupTimeout(10 * time.Second),
	}

	mgoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := mgoC.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := mgoC.MappedPort(ctx, "27017")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("mongodb://%s:%s", ip, mappedPort.Port())

	return &mongoContainer{
		Container: mgoC,
		Ip:        ip,
		Port:      mappedPort.Port(),
		URI:       uri,
	}, nil
}

func (suite *MongoTestSuite) TestNewClient() {
	client, err := NewClient(suite.mongoContainer.URI, "test")
	suite.NoError(err)
	suite.NotNil(client)
}

func (suite *MongoTestSuite) TestNewClientEmptyURI() {
	client, err := NewClient("", "test")
	suite.Error(err)
	suite.Nil(client)
}

func (suite *MongoTestSuite) TestNewClientEmptyDBName() {
	client, err := NewClient(suite.mongoContainer.URI, "")
	suite.Error(err)
	suite.Nil(client)
}

func (suite *MongoTestSuite) TestConnect() {
	client, err := NewClient(suite.mongoContainer.URI, "test")
	suite.NoError(err)
	suite.NotNil(client)

	err = client.mdbclient.Ping(context.Background(), nil)
	suite.NoError(err)
}

func (suite *MongoTestSuite) TestConnectInvalidURI() {
	client, err := NewClient("mongodb://localhost:27018", "test")
	suite.Error(err)
	suite.Nil(client)
}
