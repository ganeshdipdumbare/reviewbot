package app_test

import (
	"io"
	"log/slog"
	"testing"

	"backend/internal/app"
	convRepoMock "backend/internal/mocks/convrepo"
	msgNLPMock "backend/internal/mocks/msgnlp"
	productServiceMock "backend/internal/mocks/productservice"
	reviewRepoMock "backend/internal/mocks/reviewrepo"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewConversationApp(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	reviewRepo := reviewRepoMock.NewMockRepository(ctrl)
	convoRepo := convRepoMock.NewMockRepository(ctrl)
	msgNLP := msgNLPMock.NewMockMessageIntentService(ctrl)
	productService := productServiceMock.NewMockProductService(ctrl)

	appResp, err := app.NewConversationApp(logger, reviewRepo, convoRepo, msgNLP, productService)

	assert.NoError(t, err)
	assert.NotNil(t, appResp)

	// Test for nil reviewRepo
	appResp, err = app.NewConversationApp(logger, nil, convoRepo, msgNLP, productService)
	assert.Error(t, err)
	assert.Nil(t, appResp)

	// Test for nil convoRepo
	appResp, err = app.NewConversationApp(logger, reviewRepo, nil, msgNLP, productService)
	assert.Error(t, err)
	assert.Nil(t, appResp)

	// Test for nil msgNLP
	appResp, err = app.NewConversationApp(logger, reviewRepo, convoRepo, nil, productService)
	assert.Error(t, err)
	assert.Nil(t, appResp)

	// Test for nil productService
	appResp, err = app.NewConversationApp(logger, reviewRepo, convoRepo, msgNLP, nil)
	assert.Error(t, err)
	assert.Nil(t, appResp)

}
