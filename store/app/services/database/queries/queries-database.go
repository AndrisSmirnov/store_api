package queries

import (
	"errors"
	"fmt"

	dto "store/app/api/dto"
	types "store/app/api/types"
	postgre "store/app/services/database/connect"
	model "store/app/services/database/model"
	voc "store/app/vocabulary"

	"gorm.io/gorm"
)

func CountLengthTable(table string, count *int64) {
	postgre.ConncetDataBase.Table(table).Count(count)
}

func DeleteTable(name string) {
	postgre.ConncetDataBase.Exec("DROP TABLE " + name + " CASCADE;")
}

func UpdateTable(tx *gorm.DB, model interface{}, where, update map[string]interface{}) error {
	return tx.Model(model).Where(where).Updates(update).Error
}

func CreateUser(user *dto.ResponseUser) error {
	return postgre.ConncetDataBase.Create(user).Error
}

func SelectUser(where map[string]interface{}) *dto.ResponseUser {
	table := "users"
	selects := "*"
	user := &dto.ResponseUser{}
	postgre.ConncetDataBase.Table(table).Select(selects).Where(where).Scan(&user)

	if user.Status == voc.USER_STATUS_DELETED {
		return &dto.ResponseUser{}
	}

	return user
}

func GetUsers(where string) (userData []dto.ResponseUser) {
	postgre.ConncetDataBase.Table("users").Where(where).Scan(&userData)
	return userData
}

func CreateCategory(category *dto.ResponseCategory) error {
	return postgre.ConncetDataBase.Create(category).Error
}

func SelectCategory(where map[string]interface{}) *dto.ResponseCategory {
	table := "categories"
	selects := "*"
	category := &dto.ResponseCategory{}
	postgre.ConncetDataBase.Table(table).Select(selects).Where(where).Scan(&category)

	if category.Status == voc.CATEGORY_STATUS_DELETED {
		return &dto.ResponseCategory{}
	}

	return category
}

func CreateProduct(product *dto.ResponseProduct) error {
	return postgre.ConncetDataBase.Create(product).Error
}

func SelectProduct(where map[string]interface{}) *dto.ResponseProduct {
	table := "products"
	selects := "*"
	product := &dto.ResponseProduct{}
	postgre.ConncetDataBase.Table(table).Select(selects).Where(where).Scan(&product)

	if product.Status == voc.PRODUCT_STATUS_DELETED {
		return &dto.ResponseProduct{}
	}

	return product
}

func GetProducts() (*dto.Products, error) {
	productsTable := &[]dto.ResponseProduct{}
	table := "products"
	selects := []string{"*"}

	postgre.ConncetDataBase.Table(table).Select(selects).Scan(productsTable)

	if len(*productsTable) == 0 {
		return &dto.Products{Product: *productsTable}, errors.New(voc.HTTP_PRODUCTS_NOT_AVAILABLE)
	} else {
		return &dto.Products{Product: *productsTable}, nil
	}
}

func CreateTransaction(tx *gorm.DB, transaction *model.Transaction) (err error) {
	return tx.Model(&model.Transaction{}).Create(transaction).Error
}

func SelectTransaction(where map[string]interface{}) *dto.ResponseTransaction {
	table := "transactions"
	selects := "*"
	transaction := &dto.ResponseTransaction{}
	postgre.ConncetDataBase.Table(table).Select(selects).Where(where).Scan(&transaction)

	if transaction.Status == voc.TRANSACTION_STATUS_DELETED {
		return &dto.ResponseTransaction{}
	}

	return transaction
}

func GetTransactionLimit(selects []string, userId uint32, types types.TransactionType, time string, limit, offset int) (usersTransaction []dto.ResponseTransaction) {
	where := "user_id = ? AND types = ? AND createtime >= ?"
	postgre.ConncetDataBase.Table("transactions").Select(selects).Where(where, userId, types, time).Limit(limit).Offset(limit * offset).Scan(&usersTransaction)

	return usersTransaction
}

func FindTransaction(tx *gorm.DB, transaction *model.Transaction) (currentTransaction *dto.ResponseTransaction) {
	tx.Find(&model.Transaction{}, transaction).First(&currentTransaction)

	return currentTransaction
}

func GetLastTransaction(categories, hours int) *[]dto.WinnersTable {
	lastTransaction := &[]dto.WinnersTable{}

	table := "transactions AS t"
	selects := []string{"t.id", "t.user_id", "t.amount", "t.basket", "p.id as p_id", "p.price", "p.category_id as p_category_id"}
	joinTransaction_Product := "INNER JOIN transaction_products AS tp ON tp.transaction_id = t.id"
	joinProducts := "INNER JOIN products AS p ON tp.product_id  = p.id"
	where := fmt.Sprint("p.category_id = ? AND t.createtime >= now () - INTERVAL  '", hours, " hour'")

	postgre.ConncetDataBase.Table(table).Select(selects).Joins(joinTransaction_Product).Joins(joinProducts).Where(where, categories).Scan(&lastTransaction)

	return lastTransaction
}
