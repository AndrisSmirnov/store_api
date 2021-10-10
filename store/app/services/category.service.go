package services

import (
	"errors"

	model "store/app/services/database/model"
	queries "store/app/services/database/queries"
	voc "store/app/vocabulary"
)

func UpdateCategory(update map[string]interface{}) error {
	where := map[string]interface{}{"id": update["id"]}
	delete(update, "id")

	category := queries.SelectCategory(where)

	if category.Status == 0 {
		return errors.New(voc.HTTP_CATEGORY_NOT_FOUND)
	}

	return UpdateTable(&model.Category{}, where, update)
}
