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

func CategoryController(router *gin.Engine) {
	categoryGroup := router.Group("/category")
	categoryGroup.Use(middlewares.AuthMiddleware())
	{
		categoryGroup.POST("/create", middlewares.Validator(dto.RequestCategory{}), middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			var category *dto.ResponseCategory
			data := c.MustGet("validData").(*dto.RequestCategory)
			category.Name = data.Name
			category.Status = model.STATUS_ACTIVE

			if err := queries.CreateCategory(category); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": voc.HTTP_CATEGORY_WAS_CREATED})
		})
		categoryGroup.POST("/read", middlewares.Validator(dto.RequestFromId{}), middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			var category *dto.ResponseCategory
			categoryId := c.MustGet("validData").(*dto.RequestFromId)

			if category = queries.SelectCategory(map[string]interface{}{"id": categoryId.Id}); category.Id == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": voc.HTTP_CATEGORY_NOT_FOUND})
				return
			}

			c.JSON(http.StatusOK, category)
		})
		categoryGroup.POST("/update", middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			updateData := map[string]interface{}{}

			if err := c.BindJSON(&updateData); err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			if err := services.UpdateCategory(updateData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": voc.HTTP_CATEGORY_WAS_UPDATED})
		})
		categoryGroup.POST("/delete", middlewares.Validator(dto.RequestFromId{}), middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			categoryId := c.MustGet("validData").(*dto.RequestFromId)
			categoryData := map[string]interface{}{"id": categoryId.Id, "status": 2}

			if err := services.UpdateCategory(categoryData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": voc.HTTP_CATEGORY_DELETED})
		})
	}
}
