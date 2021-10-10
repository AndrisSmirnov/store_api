package services

import (
	"errors"
	"time"

	dto "store/app/api/dto"
	types "store/app/api/types"
	encrypt "store/app/controllers/middleware/encrypt"
	model "store/app/services/database/model"
	queries "store/app/services/database/queries"
	voc "store/app/vocabulary"
)

func CreateUser(data *dto.RequesCreatetUser) (user *dto.ResponseUser, err error) {
	user.Name = data.Name
	user.Surname = data.Surname
	user.Login = data.Login
	user.Email = data.Email
	user.Password, err = encrypt.EncryptPassword(data.Password)
	user.ParentId = data.ParentId
	user.Status = model.STATUS_ACTIVE
	user.Permission = types.Client
	user.Createtime = time.Now()
	user.Updatetime = time.Now()

	return user, err
}
func UpdateUser(update map[string]interface{}) error {
	where := map[string]interface{}{"id": update["id"]}
	delete(update, "id")

	user := queries.SelectUser(where)

	if user.Status == 0 {
		return errors.New(voc.HTTP_CLIENT_NOT_FOUND)
	}

	return UpdateTable(&model.User{}, where, update)
}
