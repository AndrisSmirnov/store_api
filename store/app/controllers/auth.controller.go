package http

import (
	"net/http"

	dto "store/app/api/dto"
	middlewares "store/app/controllers/middleware"
	services "store/app/services"

	"github.com/gin-gonic/gin"
)

func AuthController(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", middlewares.Validator(dto.RequestAuthClient{}), func(c *gin.Context) {
			data := c.MustGet("validData").(*dto.RequestAuthClient)
			services.Authentication(c, data)
		})
	}

	authGroup.Use(middlewares.AuthMiddleware())
	{
		authGroup.POST("/check", func(c *gin.Context) {
			if userId, err := services.CheckUserToken(c); err == nil {
				c.IndentedJSON(http.StatusOK, userId)
			}
		})
	}
}
