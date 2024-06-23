package app_test

import (
	"io"
	"log/slog"
	"testing"

	"backend/internal/app"
	"backend/internal/conversation"
	"backend/internal/infra/msgnlp"
	convRepoMock "backend/internal/mocks/convrepo"
	msgNLPMock "backend/internal/mocks/msgnlp"
	productServiceMock "backend/internal/mocks/productservice"
	reviewRepoMock "backend/internal/mocks/reviewrepo"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestConverse(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	reviewRepo := reviewRepoMock.NewMockRepository(ctrl)
	convoRepo := convRepoMock.NewMockRepository(ctrl)
	msgNLP := msgNLPMock.NewMockMessageIntentService(ctrl)
	productService := productServiceMock.NewMockProductService(ctrl)

	appResp, err := app.NewConversationApp(logger, reviewRepo, convoRepo, msgNLP, productService)
	assert.NoError(t, err)
	assert.NotNil(t, appResp)
}

func TestConverseIntentGreet(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	reviewRepo := reviewRepoMock.NewMockRepository(ctrl)
	convoRepo := convRepoMock.NewMockRepository(ctrl)
	msgNLP := msgNLPMock.NewMockMessageIntentService(ctrl)
	productService := productServiceMock.NewMockProductService(ctrl)

	appResp, err := app.NewConversationApp(logger, reviewRepo, convoRepo, msgNLP, productService)
	assert.NoError(t, err)
	assert.NotNil(t, appResp)

	// Test for greet intent
	convoRepo.EXPECT().Upsert(gomock.Any()).Return(&conversation.Conversation{
		ID: "someid",
	}, nil).Times(2)
	msgNLP.EXPECT().PredictIntent(&msgnlp.MessageIntentRequest{
		Text: "Hi",
	}).Return(&msgnlp.MessageIntentResponse{
		Intent: msgnlp.IntentGreet,
	}, nil).Times(1)
	converseResp, err := appResp.Converse(&app.ConverseRequest{
		Text: "Hi",
	})
	assert.NoError(t, err)
	assert.NotNil(t, converseResp)
	assert.Equal(t, "Hello!", converseResp.Text)
}

func TestConverseSystemInitiateConversation(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	reviewRepo := reviewRepoMock.NewMockRepository(ctrl)
	convoRepo := convRepoMock.NewMockRepository(ctrl)
	msgNLP := msgNLPMock.NewMockMessageIntentService(ctrl)
	productService := productServiceMock.NewMockProductService(ctrl)

	appResp, err := app.NewConversationApp(logger, reviewRepo, convoRepo, msgNLP, productService)
	assert.NoError(t, err)
	assert.NotNil(t, appResp)

	// Test for system initiate conversation intent
	convoRepo.EXPECT().Upsert(gomock.Any()).Return(&conversation.Conversation{
		ID: "someid",
	}, nil).Times(2)
	msgNLP.EXPECT().PredictIntent(&msgnlp.MessageIntentRequest{
		Text: "Hi",
	}).Return(&msgnlp.MessageIntentResponse{
		Intent: msgnlp.IntentSystemInitiateConversation,
	}, nil).Times(1)
	converseResp, err := appResp.Converse(&app.ConverseRequest{
		Text: "Hi",
	})
	assert.NoError(t, err)
	assert.NotNil(t, converseResp)
	assert.Equal(t, "Hello! How can I help you today?", converseResp.Text)
}

func TestConverseIntentGoodbye(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	reviewRepo := reviewRepoMock.NewMockRepository(ctrl)
	convoRepo := convRepoMock.NewMockRepository(ctrl)
	msgNLP := msgNLPMock.NewMockMessageIntentService(ctrl)
	productService := productServiceMock.NewMockProductService(ctrl)

	appResp, err := app.NewConversationApp(logger, reviewRepo, convoRepo, msgNLP, productService)
	assert.NoError(t, err)
	assert.NotNil(t, appResp)

	// Test for goodbye intent
	convoRepo.EXPECT().Upsert(gomock.Any()).Return(&conversation.Conversation{
		ID: "someid",
	}, nil).Times(2)
	msgNLP.EXPECT().PredictIntent(&msgnlp.MessageIntentRequest{
		Text: "Bye",
	}).Return(&msgnlp.MessageIntentResponse{
		Intent: msgnlp.IntentGoodbye,
	}, nil).Times(1)
	converseResp, err := appResp.Converse(&app.ConverseRequest{
		Text: "Bye",
	})
	assert.NoError(t, err)
	assert.NotNil(t, converseResp)
	assert.Equal(t, "Goodbye! Dont hesitate to reach out if you need help.", converseResp.Text)
}
