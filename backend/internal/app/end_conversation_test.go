package app_test

import (
	"backend/internal/app"
	convRepoMock "backend/internal/mocks/convrepo"
	msgNLPMock "backend/internal/mocks/msgnlp"
	productServiceMock "backend/internal/mocks/productservice"
	reviewRepoMock "backend/internal/mocks/reviewrepo"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestEndConversation(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	reviewRepo := reviewRepoMock.NewMockRepository(ctrl)
	convoRepo := convRepoMock.NewMockRepository(ctrl)
	msgNLP := msgNLPMock.NewMockMessageIntentService(ctrl)
	productService := productServiceMock.NewMockProductService(ctrl)

	appResp, err := app.NewConversationApp(logger, reviewRepo, convoRepo, msgNLP, productService)
	assert.NoError(t, err)
	assert.NotNil(t, appResp)

	err = appResp.EndConversation(&app.EndConversationRequest{
		ConversationID: "123",
	})
	assert.NoError(t, err)
}
