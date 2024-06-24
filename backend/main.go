package main

import (
	"log/slog"
	"os"

	"backend/config"
	"backend/internal/api/rest"
	"backend/internal/app"
	"backend/internal/infra/mongodb"
	"backend/internal/infra/msgnlp/simplenlp"
	"backend/internal/infra/productservice/simpleproduct"
	_ "backend/internal/review"
)

//	@title			Reviewbot API
//	@version		1.0
//	@description	This is a simple REST API server for Reviewbot.

//	@contact.name	Ganeshdip Dumbare
//	@contact.email	ganeshdip.dumbare@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:5174
// @BasePath	/api/v1
func main() {
	// create a slog logger instance
	opts := &slog.HandlerOptions{
		AddSource: true,
	}
	switch config.Get().LogLevel {
	case "debug":
		opts.Level = slog.LevelDebug
	case "info":
		opts.Level = slog.LevelInfo
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	}
	slogJsonHandler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(slogJsonHandler)

	// create a new db instance
	mongoClient, err := mongodb.NewClient(config.Get().MongoUri, config.Get().MongoDb)
	if err != nil {
		logger.Error("failed to create new mongo client", slog.Any("error", err))
		os.Exit(1)
	}

	// create a new review repository instance
	reviewRepo, err := mongodb.NewReviewRepository(mongoClient, config.Get().ReviewCollection)
	if err != nil {
		logger.Error("failed to create new review repository", slog.Any("error", err))
		os.Exit(1)
	}

	// create a new conversation repository instance
	converseRepo, err := mongodb.NewConversationRepository(mongoClient, config.Get().ConversationCollection)
	if err != nil {
		logger.Error("failed to create new conversation repository", slog.Any("error", err))
		os.Exit(1)
	}

	// create a new message nlp instance
	msgNLP := simplenlp.NewSimpleNLP()

	// create a new product service instance
	productService := simpleproduct.NewSimpleProductService()

	// create a new conversation app instance
	converseApp, err := app.NewConversationApp(logger, reviewRepo, converseRepo, msgNLP, productService)
	if err != nil {
		logger.Error("failed to create new conversation app", slog.Any("error", err))
		os.Exit(1)
	}

	// create a new rest api instance
	api, err := rest.NewApi(converseApp, config.Get().Port)
	if err != nil {
		logger.Error("failed to create new rest api", slog.Any("error", err))
		os.Exit(1)
	}

	// start the rest server
	api.StartServer()
}
