package http

import (
	"github.com/gin-gonic/gin"
)

func InitControllers() *gin.Engine {
	router := gin.Default()

	AuthController(router)
	UserController(router)
	CategoryController(router)
	ProductController(router)
	TransactionsController(router)

	return router
}
