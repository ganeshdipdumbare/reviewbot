package rest

import (
	appMock "backend/internal/mocks/app"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestEndConverseHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	validate = validator.New()
	appMocked := appMock.NewMockConversationApp(ctrl)
	appMocked.EXPECT().EndConversation(gomock.Any()).Return(nil)
	// Create a test router
	r := gin.Default()
	api := &apiDetails{
		app: appMocked,
	}
	r.POST("/endconverse", api.endconverse)

	// Create a request to send to the router
	jsonPayload := `{"conversationID":"123"}`
	req, _ := http.NewRequest(http.MethodPost, "/endconverse", bytes.NewBufferString(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the response
	w := httptest.NewRecorder()

	// Serve the HTTP request
	r.ServeHTTP(w, req)

	// Assert the status code and response body
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"conversationID":"123","reviewID":""}`, w.Body.String())
}
