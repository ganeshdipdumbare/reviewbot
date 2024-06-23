package rest

import (
	"log/slog"
	"net/http"

	"backend/internal/app"

	"github.com/gin-gonic/gin"
)

type endConverseRequest struct {
	ConversationID string `json:"conversationID"`
	ReviewID       string `json:"reviewID"`
}

type endConverseResponse struct {
	ConversationID string `json:"conversationID"`
	ReviewID       string `json:"reviewID"`
}

// endconverse godoc
// @Summary end conversation
// @Description end conversation with the bot
// @Tags conversation
// @Accept  json
// @Produce  json
// @Param endConverseRequest body endConverseRequest true "end conversation with the bot"
// @Success 201 {object} converseResponse
// @Failure 400 {object} errorRespose
// @Failure 500 {object} errorRespose
// @Router /endconverse [post]
func (api *apiDetails) endconverse(c *gin.Context) {
	req := &endConverseRequest{}
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

	err = api.app.EndConversation(&app.EndConversationRequest{
		ConversationID: req.ConversationID,
		ReviewID:       req.ReviewID,
	})
	if err != nil {
		slog.Debug("error in getting response from app", err)
		createErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusCreated, &endConverseResponse{
		ConversationID: req.ConversationID,
		ReviewID:       req.ReviewID,
	})
	c.Done()
}
