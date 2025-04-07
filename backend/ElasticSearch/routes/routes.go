package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/saveES", savetoES)
	server.GET("/search-http", searchES)
	server.GET("/search", func(ctx *gin.Context) {
		searchESWS(ctx)
	})

}
