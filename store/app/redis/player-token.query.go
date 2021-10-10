package redis

import (
	"fmt"
	"strconv"
)

func (session *RedisSession) SaveUserToken(Id uint32, accessToken string) error {
	err := session.Client.Set(
		session.Ctx,
		fmt.Sprint("user-valid-token:", strconv.Itoa(int(Id))),
		accessToken,
		0,
	).Err()

	return err
}

func (session *RedisSession) GetTokenFromUserId(Id uint32) (string, error) {
	accessToken, err := session.Client.Get(
		session.Ctx,
		fmt.Sprint("user-valid-token:", strconv.Itoa(int(Id))),
	).Result()

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (session *RedisSession) DeleteTokenByUserLogin(Id uint32) {
	session.Client.Del(
		session.Ctx,
		fmt.Sprint("user-valid-token:", strconv.Itoa(int(Id))),
	)
}
