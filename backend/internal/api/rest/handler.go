package rest

import (
	_ "backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	swagFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var validate *validator.Validate

type errorRespose struct {
	ErrorMessage string `json:"errorMessage"`
}

func createErrorResponse(c *gin.Context, code int, message string) {
	c.IndentedJSON(code, &errorRespose{
		ErrorMessage: message,
	})
}

func (api *apiDetails) setupRouter() *gin.Engine {
	validate = validator.New()
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))
	apiV1 := r.Group("/api/v1")
	apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swagFiles.Handler))
	apiV1.POST("/converse", api.converse)
	apiV1.POST("/endconverse", api.endconverse)
	apiV1.GET("/health", api.health)
	return r
}
