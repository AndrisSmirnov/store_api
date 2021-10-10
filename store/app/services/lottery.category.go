package services

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	dto "store/app/api/dto"
	types "store/app/api/types"
	postgre "store/app/services/database/connect"
	model "store/app/services/database/model"
	queries "store/app/services/database/queries"
)

func GetWinners(procent, maxWinners, hours int) {
	lastTransaction := GetLastTransaction(RandomCategory(), hours)
	winnersMapUserAmount := CreateUserAmountMap(lastTransaction)

	if len(winnersMapUserAmount) == 0 {
		return
	}

	prize := CalculatePrize(winnersMapUserAmount, procent, maxWinners)

	if err := BringPrizeToWinners(winnersMapUserAmount, prize, maxWinners); err != nil {
		GetWinners(procent, maxWinners, hours)
	}
}

func RandomCategory() (number int) {
	var count int64

	queries.CountLengthTable("categories", &count)
	rand.Seed(time.Now().UnixNano())
	number = randInt(1, int(count)+1)

	// TODO Clear all fmt
	fmt.Println("\nNUMBER OF CATEGORY:\t", number)
	return
}

func GetLastTransaction(category, hours int) []dto.WinnersTable {
	return *queries.GetLastTransaction(category, hours)
}

func CreateUserAmountMap(lastTransactions []dto.WinnersTable) map[string]int {
	var basket []map[string]interface{}
	winnersMapUserAmount := map[string]int{}

	for i := 0; i != len(lastTransactions); i++ {
		byt := []byte(lastTransactions[i].Basket)

		if err := json.Unmarshal(byt, &basket); err != nil {
			log.Println(err.Error())
		}

		count := FoundCountOrders(basket, lastTransactions[i].PId)
		winnersMapUserAmount[fmt.Sprint(lastTransactions[i].UserId)] += count * int(lastTransactions[i].Price)
	}

	return winnersMapUserAmount
}

// TODO Clean all fmt
func CalculatePrize(winners map[string]int, procent, maxWinners int) (prize int) {
	currentCountWinners := 0

	for _, value := range winners {
		currentCountWinners++
		prize += value
	}

	fmt.Println("Max Winners:\t\t", maxWinners, "\nCurrentCountWinners:\t", currentCountWinners)
	if currentCountWinners >= maxWinners {
		fmt.Println("\nSum of Category:\t", prize, "\nPrize/Procent:\t\t", prize/(procent*maxWinners))
		return prize / (procent * maxWinners)
	} else {
		fmt.Println("\nSum of Category:\t", prize, "\nPrize/Procent:\t\t", prize/(procent*currentCountWinners))
		return prize / (procent * currentCountWinners)
	}
}

func BringPrizeToWinners(winnersMap map[string]int, prize, maxWinners int) (err error) {
	return BringPrize(winnersMap, prize, maxWinners)
}

func FoundCountOrders(basket []map[string]interface{}, produstId uint32) (count int) {
	for i := 0; i != len(basket); i++ {
		if uint32(basket[i]["productId"].(float64)) == produstId {
			count = int(basket[i]["count"].(float64))
			return count
		}
	}

	return count
}

func BringPrize(winnersMap map[string]int, prize, maxWinners int) (err error) {
	where, updateRaw := map[string]interface{}{}, map[string]interface{}{}
	var userBalance uint32
	var basket string
	var createTime time.Time

	winners := FindWinners(winnersMap, maxWinners)
	// TODO Clean all fmt
	fmt.Println("\nWINNERS:\t\t", winners)

	requestUsers := "id IN ("
	for i := 0; i != len(winners)-1; i++ {
		requestUsers = fmt.Sprint(requestUsers, winners[i].User, ", ")
	}
	requestUsers = fmt.Sprint(requestUsers, winners[len(winners)-1].User, ")")

	winnersData := queries.GetUsers(requestUsers)

	tx := postgre.ConncetDataBase.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for i := 0; i != len(winners); i++ {
		createTime = time.Now()
		basket = fmt.Sprint("Winner ", winnersData[i].Id, " with Prize: ", prize, " at ", createTime)
		userBalance = uint32(winnersData[i].Balance + uint32(prize))
		newTransaction := &model.Transaction{
			Balance:    userBalance,
			Amount:     uint32(prize),
			Types:      uint32(types.TransactionType_WIN),
			Status:     model.STATUS_ACTIVE,
			UserId:     winnersData[i].Id,
			Basket:     basket,
			Createtime: createTime,
			Updatetime: time.Now(),
		}

		if err := tx.Model(&model.Transaction{}).Create(newTransaction).Error; err != nil {
			tx.Rollback()
			return err
		}

		where["id"] = winnersData[i].Id
		updateRaw["balance"] = userBalance
		if err := tx.Model(&model.User{}).Where(where).Updates(updateRaw).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func FindWinners(winnersMap map[string]int, maxWinners int) (winnersSortedStruct []types.LotteryUserAmount) {
	sorted_struct := []types.LotteryUserAmount{}

	for key, value := range winnersMap {
		sorted_struct = append(sorted_struct, types.LotteryUserAmount{User: key, Amount: value})
	}

	sort.Slice(sorted_struct, func(i, j int) bool {
		return sorted_struct[i].Amount > sorted_struct[j].Amount
	})

	if len(winnersMap) < maxWinners {
		maxWinners = len(winnersMap)
	}

	for i := 0; i != maxWinners; i++ {
		winnersSortedStruct = append(winnersSortedStruct, sorted_struct[i])
	}

	return winnersSortedStruct
}
