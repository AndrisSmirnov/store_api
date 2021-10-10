package services

import (
	"errors"

	dto "store/app/api/dto"
	model "store/app/services/database/model"
	queries "store/app/services/database/queries"
	voc "store/app/vocabulary"
)

func UpdateProduct(update map[string]interface{}) error {
	where := map[string]interface{}{"id": update["id"]}
	delete(update, "id")

	product := queries.SelectProduct(where)

	if product.Status == 0 {
		return errors.New(voc.HTTP_PRODUCT_NOT_FOUND)
	}

	return UpdateTable(&model.Product{}, where, update)
}

func GetRawFromProductTable(number uint32, productTable *dto.Products) (raw *dto.ResponseProduct) {
	for i := 0; i != len(productTable.Product); i++ {
		raw = &productTable.Product[i]
		if raw.Id == number {
			break
		}
	}

	return raw
}
