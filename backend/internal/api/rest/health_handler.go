package rest

import "github.com/gin-gonic/gin"

func (api *apiDetails) health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
