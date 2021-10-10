package services

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	dto "store/app/api/dto"
	types "store/app/api/types"
	encrypt "store/app/controllers/middleware/encrypt"
	redis "store/app/redis"
	queries "store/app/services/database/queries"
	voc "store/app/vocabulary"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context, auth *dto.RequestAuthClient) (string, error) {
	client, err := GetClient(*auth)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return "", err
	}

	if prevToken, err := CheckIfTokenExistRedis(client.Id); err == nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": voc.REDIS_TOKEN_EXIST})
		return prevToken, nil
	}

	var expiresKey string

	if client.Permission >= types.Admin {
		expiresKey = "ADMIN_SECRET_KEY_EXPIRES"
	} else {
		expiresKey = "TOKEN_EXPIRES"
	}

	expires, err := strconv.ParseInt(os.Getenv(expiresKey), 10, 64)

	if err != nil {
		return "", err
	}

	claims := types.Token{
		Id:         client.Id,
		Permission: client.Permission,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(expires)).UnixNano(),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	res, err := token.SignedString([]byte(os.Getenv("TOKEN_SALT")))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return "", err
	}

	if err = redis.Session.SaveUserToken(claims.Id, res); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return "", err
	}

	c.IndentedJSON(http.StatusOK, res)
	return res, nil
}

func CheckToken(tokenString string) (*types.Token, error) {
	res := types.Token{}

	token, err := jwt.ParseWithClaims(tokenString, &res, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SALT")), nil
	})

	claims, ok := token.Claims.(*types.Token)

	if !ok || !token.Valid {
		return nil, err
	}

	if time.Now().UnixNano() > claims.ExpiresAt {
		return nil, errors.New(voc.TOKEN_EXPIRED)
	}

	return &res, nil
}

func CheckUserToken(c *gin.Context) (string, error) {
	token := c.MustGet("jwtPayloadData").(*types.Token)

	if isStatusDeleted := CheckUserStatus(token.Id); isStatusDeleted {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": voc.HTTP_CLIENT_DELETED})
		return "", errors.New(voc.HTTP_CLIENT_DELETED)
	}

	if _, err := redis.Session.GetTokenFromUserId(token.Id); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": voc.ERROR_REDIS_NIL_OBJ})
		return "", errors.New(voc.ERROR_REDIS_NIL_OBJ)
	}

	idString := strconv.Itoa(int(token.Id))
	return idString, nil
}

func CheckUserStatus(userId uint32) bool {
	where := map[string]interface{}{"id": userId}
	user := queries.SelectUser(where)

	if user.Password == "" {
		PrintLogsSomeoneHasSalt(userId)
		return true
	}

	if user.Status == voc.USER_STATUS_DELETED {
		return true
	} else {
		return false
	}
}

// TODO: model.User or dto.User for response from DB?
func GetClient(data dto.RequestAuthClient) (user *dto.ResponseUser, err error) {
	where := map[string]interface{}{"login": data.Login}
	user = queries.SelectUser(where)

	if user.Login == "" {
		return user, errors.New(voc.HTTP_CLIENT_NOT_FOUND)
	}

	if !encrypt.CheckPasswordHash(data.Password, user.Password) {
		return user, errors.New(voc.HTTP_WROND_PASSWORD)
	}

	if user.Status == voc.USER_STATUS_DELETED {
		return user, errors.New(voc.HTTP_CLIENT_DELETED)
	}

	return user, nil
}

func CheckIfTokenExistRedis(id uint32) (string, error) {
	prevToken, err := redis.Session.GetTokenFromUserId(id)

	if err != nil {
		return "", err
	}

	if _, err = CheckToken(prevToken); err != nil {
		redis.Session.DeleteTokenByUserLogin(id)
		return "", err
	}

	return prevToken, nil
}
