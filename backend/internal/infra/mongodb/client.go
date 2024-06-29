package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoConnectTimeout = 10 * time.Second
)

var ErrEmptyArg = errors.New("empty argument")

type Client struct {
	mdbclient *mongo.Client
	dbName    string
}

// NewClient created new mongo db instance, returns error if input is invalid
func NewClient(uri string, dbName string) (*Client, error) {
	if uri == "" {
		return nil, fmt.Errorf("%w: uri", ErrEmptyArg)
	}

	if dbName == "" {
		return nil, fmt.Errorf("%w: dbName", ErrEmptyArg)
	}

	mdbclient, err := connect(uri)
	if err != nil {
		return nil, err
	}

	return &Client{
		mdbclient: mdbclient,
		dbName:    dbName,
	}, nil
}

// connect connects to mongo db using client, returns error if fails
func connect(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), mongoConnectTimeout)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
