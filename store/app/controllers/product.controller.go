package http

import (
	"net/http"

	dto "store/app/api/dto"
	types "store/app/api/types"
	middlewares "store/app/controllers/middleware"
	services "store/app/services"
	model "store/app/services/database/model"
	queries "store/app/services/database/queries"
	voc "store/app/vocabulary"

	"github.com/gin-gonic/gin"
)

func ProductController(router *gin.Engine) {
	productGroup := router.Group("/product")
	productGroup.Use(middlewares.AuthMiddleware())
	{
		productGroup.POST("/create", middlewares.Validator(dto.RequestProduct{}), middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			var product *dto.ResponseProduct
			data := c.MustGet("validData").(*dto.RequestProduct)
			product.Name = data.Name
			product.Price = data.Price
			product.Count = data.Count
			product.Status = model.STATUS_ACTIVE
			product.CategoryId = data.CategoryId
			if err := queries.CreateProduct(product); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": voc.HTTP_PRODUCT_WAS_CREATED})
		})
		productGroup.POST("/read", middlewares.Validator(dto.RequestFromId{}), middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			var product *dto.ResponseProduct
			productId := c.MustGet("validData").(*dto.RequestFromId)

			if product = queries.SelectProduct(map[string]interface{}{"id": productId.Id}); product.Id == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": voc.HTTP_PRODUCT_NOT_FOUND})
				return
			}

			c.JSON(http.StatusOK, product)
		})
		productGroup.POST("/update", middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			updateData := map[string]interface{}{}

			if err := c.BindJSON(&updateData); err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			if err := services.UpdateProduct(updateData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": voc.HTTP_PRODUCT_WAS_UPDATED})
		})
		productGroup.POST("/delete", middlewares.Validator(dto.RequestFromId{}), middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			productId := c.MustGet("validData").(*dto.RequestFromId)
			productData := map[string]interface{}{"id": productId.Id, "status": 2}

			if err := services.UpdateProduct(productData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": voc.HTTP_PRODUCT_DELETED})
		})
	}
}
