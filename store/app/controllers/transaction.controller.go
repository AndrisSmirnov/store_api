package http

import (
	"net/http"

	dto "store/app/api/dto"
	types "store/app/api/types"
	middlewares "store/app/controllers/middleware"
	services "store/app/services"
	queries "store/app/services/database/queries"
	voc "store/app/vocabulary"

	"github.com/gin-gonic/gin"
)

func TransactionsController(router *gin.Engine) {
	transactionGroup := router.Group("/transaction")
	transactionGroup.Use(middlewares.AuthMiddleware())
	{
		transactionGroup.POST("/bye", middlewares.Validator(dto.RequestTransaction{}), middlewares.PermissionMiddleware(types.Client), func(c *gin.Context) {
			order := c.MustGet("validData").(*dto.RequestTransaction)

			if idUser, err := services.CheckUserToken(c); err == nil {
				services.CreateTransactionDB(c, order, idUser)
			}
		})
		transactionGroup.POST("/gettransaction", middlewares.Validator(dto.RequestShowTransaction{}), middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			data := c.MustGet("validData").(*dto.RequestShowTransaction)

			services.GetTransactionLimitDB(data, c)
		})
		transactionGroup.POST("/read", middlewares.Validator(dto.RequestFromId{}), middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			var transaction *dto.ResponseTransaction
			transactionId := c.MustGet("validData").(*dto.RequestFromId)

			if transaction = queries.SelectTransaction(map[string]interface{}{"id": transactionId.Id}); transaction.Id == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": voc.HTTP_TRANSACTION_NOT_FOUND})
				return
			}

			c.JSON(http.StatusOK, transaction)
		})
		transactionGroup.POST("/update", middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			updateData := map[string]interface{}{}

			if err := c.BindJSON(&updateData); err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			if err := services.UpdateTransaction(updateData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": voc.HTTP_CATEGORY_WAS_UPDATED})
		})
		transactionGroup.POST("/delete", middlewares.Validator(dto.RequestFromId{}), middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			transactionId := c.MustGet("validData").(*dto.RequestFromId)
			transactionData := map[string]interface{}{"id": transactionId.Id, "status": 2}

			if err := services.UpdateTransaction(transactionData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": voc.HTTP_TRANSACTION_DELETED})
		})
	}

}
