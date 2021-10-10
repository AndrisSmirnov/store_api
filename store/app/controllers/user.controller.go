package http

import (
	"net/http"

	dto "store/app/api/dto"
	types "store/app/api/types"
	middlewares "store/app/controllers/middleware"
	services "store/app/services"
	fill_db "store/app/services/database/filldatabase"
	queries "store/app/services/database/queries"
	voc "store/app/vocabulary"

	"github.com/gin-gonic/gin"
)

func UserController(router *gin.Engine) {
	userGroup := router.Group("/user")
	userGroup.Use(middlewares.AuthMiddleware())
	{
		userGroup.POST("/create", middlewares.Validator(dto.RequesCreatetUser{}), middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			data := c.MustGet("validData").(*dto.RequesCreatetUser)
			user, err := services.CreateUser(data)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "failed"})
				return
			}

			if err = queries.CreateUser(user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": voc.HTTP_USER_WAS_CREATED})
		})
		userGroup.POST("/read", middlewares.Validator(dto.RequestFromId{}), middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			var user *dto.ResponseUser
			userId := c.MustGet("validData").(*dto.RequestFromId)

			if user = queries.SelectUser(map[string]interface{}{"id": userId.Id}); user.Id == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": voc.HTTP_CLIENT_NOT_FOUND})
				return
			}

			c.JSON(http.StatusOK, user)
		})
		userGroup.POST("/update", middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			updateData := map[string]interface{}{}

			if err := c.BindJSON(&updateData); err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			if err := services.UpdateUser(updateData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": voc.HTTP_USER_WAS_UPDATED})
		})
		userGroup.POST("/delete", middlewares.Validator(dto.RequestFromId{}), middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			userId := c.MustGet("validData").(*dto.RequestFromId)
			updateData := map[string]interface{}{"id": userId.Id, "status": 2}

			if err := services.UpdateUser(updateData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": voc.HTTP_CLIENT_DELETED})
		})
		userGroup.POST("/checkparents", middlewares.Validator(dto.RequestParents{}), middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			data := c.MustGet("validData").(*dto.RequestParents)

			parents, _ := fill_db.GetParentsStack(data.UserId)
			c.IndentedJSON(http.StatusOK, parents)
		})
		userGroup.POST("/checkchilndrens", middlewares.Validator(dto.RequestChildrens{}), middlewares.PermissionMiddleware(types.Admin), func(c *gin.Context) {
			data := c.MustGet("validData").(*dto.RequestChildrens)

			childrens, childrensbalance := fill_db.CheckChilndrensStackAndBallance(data.UserId, int(data.Deep))

			switch val := int(data.Type); val {
			case 1:
				c.IndentedJSON(http.StatusOK, gin.H{"Childrens": childrens})
			case 2:
				c.IndentedJSON(http.StatusOK, gin.H{"ChildrensBalance": childrensbalance})
			case 3:
				childrensSUM := services.SumOf2DArray(childrensbalance)
				c.IndentedJSON(http.StatusOK, gin.H{"ChildrensBalanceSUM": childrensSUM})
			default:
				c.IndentedJSON(http.StatusBadGateway, voc.HTTP_INVALID_REQUEST)
			}
		})
	}
}
