/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"

	"backend/config"
	"backend/internal/api/rest"
	"backend/internal/app"
	"backend/internal/infra/mongodb"
	"backend/internal/infra/msgnlp/simplenlp"
	"backend/internal/infra/productservice/simpleproduct"

	"github.com/spf13/cobra"
)

var port int

// restCmd represents the rest command
var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "A command to start rest api server",
	Long: `A command to start rest api server at given port. 
	If the port is not provided, it will start the server at default port 8080.`,
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}

func init() {
	rootCmd.AddCommand(restCmd)

	restCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
}
