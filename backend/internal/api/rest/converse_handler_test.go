package rest

import (
	"backend/internal/app"
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

func TestConverseHandlerGreet(t *testing.T) {
	ctrl := gomock.NewController(t)
	validate = validator.New()
	appMocked := appMock.NewMockConversationApp(ctrl)
	appMocked.EXPECT().Converse(gomock.Any()).Return(&app.ConverseResponse{
		ConversationID: "123",
		ReviewID:       "456",
		UserID:         "789",
		ProductID:      "101112",
		Text:           "Hello!",
	}, nil)
	// Create a test router
	r := gin.Default()
	api := &apiDetails{
		app: appMocked,
	}
	r.POST("/converse", api.converse)

	// Create a request to send to the router
	jsonPayload := `{"text":"Hello"}`
	req, _ := http.NewRequest(http.MethodPost, "/converse", bytes.NewBufferString(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the response
	w := httptest.NewRecorder()

	// Serve the HTTP request
	r.ServeHTTP(w, req)

	// Assert the status code and response body
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t,
		`{
    "conversationID": "123",
    "reviewID": "456",
    "userID": "789",
    "productID": "101112",
    "text": "Hello!"
	}`, w.Body.String())
}
