package middleware

import (
	"net/http"

	dto "store/app/api/dto"
	model "store/app/services/database/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
)

func Validator(i interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data interface{}
		switch i.(type) {
		case dto.RequestFromId:
			data = &dto.RequestFromId{}
		//auth
		case dto.RequestAuthClient:
			data = &dto.RequestAuthClient{}
		//user
		case dto.RequesCreatetUser:
			data = &dto.RequesCreatetUser{}
		case dto.RequestParents:
			data = &dto.RequestParents{}
		case dto.RequestChildrens:
			data = &dto.RequestChildrens{}
		//category
		case dto.RequestCategory:
			data = &dto.RequestCategory{}
		//product
		case dto.RequestProduct:
			data = &dto.RequestProduct{}
		//transaction
		case dto.RequestTransaction:
			data = &dto.RequestTransaction{}
		case dto.ResponseTransaction:
			data = &model.Transaction{}
		case dto.RequestShowTransaction:
			data = &dto.RequestShowTransaction{}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "dto type is invalid"})
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := c.Bind(data); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		v := validator.New()

		_ = v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
			return len(fl.Field().String()) > 6
		})

		if err := v.Struct(data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		conform := modifiers.New()

		if err := conform.Struct(c, data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Set("validData", data)
		c.Next()
	}
}
