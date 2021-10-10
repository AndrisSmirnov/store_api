package middleware

import (
	"errors"
	"net/http"
	"strings"

	types "store/app/api/types"
	services "store/app/services"
	voc "store/app/vocabulary"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := checkToken(c.Request.Header.Get("Authorization"))

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("jwtPayloadData", user)

		c.Next()
	}
}

func PermissionMiddleware(level types.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := c.MustGet("jwtPayloadData").(*types.Token)

		if data.Permission < level {
			c.JSON(http.StatusForbidden, gin.H{"error": voc.PERMISSION_ERROR})
			c.Abort()
			return
		}
	}

}

func checkToken(auth string) (*types.Token, error) {
	if auth == "" {
		return nil, errors.New(voc.HTTP_NO_HEADER)
	}

	tokenString := strings.TrimPrefix(auth, "Bearer ")

	if tokenString == auth {
		return nil, errors.New(voc.HTTP_WRONG_HEADER)
	}

	user, err := services.CheckToken(tokenString)

	if err != nil {
		return nil, err
	}

	return user, nil
}
