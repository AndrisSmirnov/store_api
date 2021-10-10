package services

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	postgre "store/app/services/database/connect"
	queries "store/app/services/database/queries"

	"time"
)

func UpdateTable(models interface{}, where, update map[string]interface{}) error {
	tx := postgre.ConncetDataBase.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := queries.UpdateTable(tx, models, where, update); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func SumOf2DArray(arr [][]uint32) (sum uint32) {
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			sum += arr[i][j]
		}
	}

	return sum
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func PrintLogsSomeoneHasSalt(userId uint32) {
	logs := fmt.Sprint("\nsomeone can generate a token without registration\nid: \t", userId, "\n\n ")
	log.Println(logs)

	f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(fmt.Sprint(time.Now(), logs)); err != nil {
		panic(err)
	}
}
