package services

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	dto "store/app/api/dto"
	types "store/app/api/types"
	postgre "store/app/services/database/connect"
	model "store/app/services/database/model"
	queries "store/app/services/database/queries"
	voc "store/app/vocabulary"

	"github.com/gin-gonic/gin"
)

func UpdateTransaction(update map[string]interface{}) error {
	where := map[string]interface{}{"id": update["id"]}
	delete(update, "id")

	transaction := queries.SelectTransaction(where)

	if transaction.Status == 0 {
		return errors.New(voc.HTTP_TRANSACTION_NOT_FOUND)
	}

	return UpdateTable(&model.Transaction{}, where, update)
}

func GetTransactionLimitDB(requestTrans *dto.RequestShowTransaction, c *gin.Context) {
	location, _ := time.LoadLocation("UTC")
	tm := time.Unix(int64(requestTrans.UnixTime), 0).In(location).Format("2006-01-02 15:04:05")
	transactions := queries.GetTransactionLimit([]string{"*"}, requestTrans.UserId, requestTrans.Type, tm, int(requestTrans.Limit), int(requestTrans.Offset-1))

	if len(transactions) != 0 {
		c.JSON(http.StatusOK, gin.H{"transactions": transactions})
	}
}

func CreateTransactionDB(c *gin.Context, order *dto.RequestTransaction, userId string) {
	var createTime time.Time = time.Now()

	productsTable, err := queries.GetProducts()

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	userData := queries.SelectUser(map[string]interface{}{"id": userId})

	switch result := false; result {
	case CheckProductsInStock(c, productsTable, order):
		return
	case CheckUserSolvency(c, productsTable, order, userData):
		return
	case Trascaction(c, productsTable, order, userData, createTime):
		return
	default:
		c.IndentedJSON(http.StatusBadRequest, voc.HTTP_TRANSACTION_OK)
	}
}

func CheckProductsInStock(c *gin.Context, productsTable *dto.Products, order *dto.RequestTransaction) (succsess bool) {
	var err error = nil

	for i := 0; i != len(order.RequestOrder); i++ {
		productRaw := GetRawFromProductTable(order.RequestOrder[i].ProductId, productsTable)

		if productRaw == nil {
			err = errors.New(voc.HTTP_INVALID_REQUEST)
			break
		}

		if productRaw.Count < order.RequestOrder[i].Count {
			err = errors.New(voc.HTTP_PRODUCTS_OUT_STOCK)
			break
		}

		if productRaw.Status == voc.PRODUCT_STATUS_DELETED {
			err = errors.New(voc.HTTP_PRODUCTS_NOT_AVAILABLE)
			break
		}
	}

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return succsess
	}

	succsess = true
	return succsess
}

func CheckUserSolvency(c *gin.Context, productsTable *dto.Products, order *dto.RequestTransaction, userData *dto.ResponseUser) (succsess bool) {
	var err error = nil
	var sum uint32 = 0

	for i := 0; i != len(order.RequestOrder); i++ {
		product := GetRawFromProductTable(order.RequestOrder[i].ProductId, productsTable)

		sum += product.Price * order.RequestOrder[i].Count
	}

	if userData.Balance < sum {
		err = errors.New(voc.HTTP_CLIENT_WITHOUT_FOUNDS)
	}

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return succsess
	}

	succsess = true
	return succsess
}

func Trascaction(c *gin.Context, productsTable *dto.Products, order *dto.RequestTransaction, userData *dto.ResponseUser, createTime time.Time) (succsess bool) {
	if err := CreateTransaction(productsTable, order, userData, createTime); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return succsess
	}

	succsess = true
	return succsess
}

func CreateTransaction(productsTable *dto.Products, order *dto.RequestTransaction, userData *dto.ResponseUser, createTime time.Time) (err error) {
	var amount uint32 = 0
	var where, updateRaw map[string]interface{} = map[string]interface{}{"id": userData.Id}, map[string]interface{}{}
	transactionProdacts := make([]dto.TransactionProdacts, len(order.RequestOrder))

	basket := CreateBasket(order)

	CountAmount(productsTable, order, &amount)

	tx := postgre.ConncetDataBase.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	newTransaction := &model.Transaction{
		Balance:    userData.Balance - amount,
		Amount:     amount,
		Types:      uint32(types.TransactionType_DEC),
		Status:     model.STATUS_ACTIVE,
		UserId:     userData.Id,
		Basket:     basket,
		Createtime: createTime,
		Updatetime: time.Now(),
	}

	if err := queries.CreateTransaction(tx, newTransaction); err != nil {
		tx.Rollback()
		return err
	}

	updateRaw["balance"] = userData.Balance - amount

	if err := queries.UpdateTable(tx, &model.User{}, where, updateRaw); err != nil {
		tx.Rollback()
		return err
	}

	currentTransaction := queries.FindTransaction(tx, newTransaction)

	for i := range order.RequestOrder {
		transactionProdacts[i].TransactionId = currentTransaction.Id
		transactionProdacts[i].ProductId = order.RequestOrder[i].ProductId
		where["id"] = order.RequestOrder[i].ProductId
		updateRaw = GetUpdateProductsRaw(productsTable, order.RequestOrder[i])

		if err := queries.UpdateTable(tx, &model.Product{}, where, updateRaw); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Table("transaction_products").Create(&transactionProdacts).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func GetUpdateProductsRaw(productsTable *dto.Products, order dto.RequestOrder) map[string]interface{} {
	updateRaw := map[string]interface{}{}

	for i := range productsTable.Product {
		if order.ProductId == productsTable.Product[i].Id {
			updateRaw["count"] = productsTable.Product[i].Count - order.Count
			return updateRaw
		}
	}

	return updateRaw
}

func CreateBasket(order *dto.RequestTransaction) (basket string) {
	basket = "["

	for i := range order.RequestOrder {
		basket = fmt.Sprint(basket, "{\"productId\":", order.RequestOrder[i].ProductId, ", \"count\":", order.RequestOrder[i].Count, "}")
		if i != len(order.RequestOrder)-1 {
			basket = fmt.Sprint(basket, ",")
		}
	}

	basket = fmt.Sprint(basket, "]")

	return basket
}

func CountAmount(productsTable *dto.Products, order *dto.RequestTransaction, amount *uint32) {
	var price uint32

	for i := range order.RequestOrder {
		for j := range productsTable.Product {
			if order.RequestOrder[i].ProductId == productsTable.Product[j].Id {
				price = productsTable.Product[j].Price
				break
			}
		}
		*amount += order.RequestOrder[i].Count * price
	}
}
