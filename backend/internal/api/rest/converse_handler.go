package rest

import (
	"log/slog"
	"net/http"

	"backend/internal/app"

	"github.com/gin-gonic/gin"
)

type converseRequest struct {
	ConversationID string `json:"conversationID"`
	ReviewID       string `json:"reviewID"`
	UserID         string `json:"userID"`
	ProductID      string `json:"productID"`
	Text           string `json:"text" binding:"required"`
}

type converseResponse struct {
	ConversationID string `json:"conversationID"`
	ReviewID       string `json:"reviewID"`
	UserID         string `json:"userID"`
	ProductID      string `json:"productID"`
	Text           string `json:"text"`
}

// converse godoc
// @Summary add message to conversation and get response
// @Description add message to conversation and get response from the bot
// @Tags conversation
// @Accept  json
// @Produce  json
// @Param converseRequest body converseRequest true "continue conversation with the bot"
// @Success 201 {object} converseResponse
// @Failure 400 {object} errorRespose
// @Failure 500 {object} errorRespose
// @Router /converse [post]
func (api *apiDetails) converse(c *gin.Context) {
	req := &converseRequest{}
	err := c.BindJSON(req)
	if err != nil {
		slog.Debug("error in binding json", err)
		createErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.Struct(req)
	if err != nil {
		createErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := api.app.Converse(&app.ConverseRequest{
		Text:           req.Text,
		ConversationID: req.ConversationID,
		ReviewID:       req.ReviewID,
		UserID:         req.UserID,
		ProductID:      req.ProductID,
	})
	if err != nil {
		slog.Debug("error in getting response from app", err)
		createErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusCreated, &converseResponse{
		ConversationID: resp.ConversationID,
		ReviewID:       resp.ReviewID,
		UserID:         resp.UserID,
		ProductID:      resp.ProductID,
		Text:           resp.Text,
	})
	c.Done()
}
