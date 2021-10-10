package filldatabase

import (
	"fmt"
	"math/rand"

	"strconv"
	"time"

	dto "store/app/api/dto"
	types "store/app/api/types"
	encrypt "store/app/controllers/middleware/encrypt"
	postgre "store/app/services/database/connect"
	model "store/app/services/database/model"
	queries "store/app/services/database/queries"
	voc "store/app/vocabulary"
)

var category = []model.Category{{Name: "Dairy", Status: 1}, {Name: "Meat", Status: 1}, {Name: "Fruits", Status: 1}}
var products = []model.Product{
	{Name: "Milk", Price: 35, Count: 13, Status: 2, CategoryId: 3},
	{Name: "Cheese", Price: 200, Count: 27, Status: 2, CategoryId: 1},
	{Name: "Butter", Price: 250, Count: 14, Status: 1, CategoryId: 1},
	{Name: "Becon", Price: 150, Count: 30, Status: 1, CategoryId: 2},
	{Name: "Chicken", Price: 100, Count: 11, Status: 1, CategoryId: 2},
	{Name: "Apple", Price: 15, Count: 49, Status: 1, CategoryId: 3},
	{Name: "Melon", Price: 70, Count: 6, Status: 1, CategoryId: 3},
}

func FillDataBase(users, maxLengthTree int) {
	CreateRowsCategories(category)
	CreateRowsProducts(products)
	CreateNextUsersDB(users, maxLengthTree)
}

func CreateRowsCategories(category []model.Category) {
	postgre.ConncetDataBase.Model(&model.Category{}).Create(&category)
}

func CreateRowsProducts(products []model.Product) {
	postgre.ConncetDataBase.Model(&model.Product{}).Create(&products)
}

func CreateNextUsersDB(howMuch, maxLengthTree int) {
	var counterUsers int
	var countTableLength int64

	queries.CountLengthTable("users", &countTableLength)
	counterUsers = int(countTableLength) + 1

	if countTableLength != 0 {
		for howMuch != 0 {
			CreateUser(counterUsers, int(countTableLength), maxLengthTree, CreatePassword(counterUsers))
			howMuch -= 1
			counterUsers++
		}
	} else {
		if howMuch != 0 {
			CreateUser(1, int(countTableLength), maxLengthTree, CreatePassword(counterUsers))
			CreateNextUsersDB(howMuch-1, maxLengthTree)
		}
	}
}

func CreatePassword(counterUsers int) (pass string) {
	pass, _ = encrypt.EncryptPassword("password" + strconv.Itoa(counterUsers))
	return pass
}

func CreateUser(counterUsers, tableUserLength, maxLengthTree int, password string) {
	postgre.ConncetDataBase.Model(&model.User{}).Create(
		&model.User{
			Name:       strconv.Itoa(counterUsers) + "Name",
			Surname:    strconv.Itoa(counterUsers) + "Surname",
			Login:      strconv.Itoa(counterUsers) + "Login",
			Email:      strconv.Itoa(counterUsers) + "mail@gmail.com",
			Password:   password,
			Balance:    uint32(counterUsers) * 1000,
			Status:     voc.USER_STATUS_ACTIVE,
			Permission: types.Client,
			ParentId:   CreateParent(counterUsers, tableUserLength, maxLengthTree)},
	)
}

func CreateParent(elmentUserId, tableUserLength, maxLengthTree int) uint32 {
	var parent *dto.ResponseUser
	randomParentId := 0
	where := map[string]interface{}{}

	if tableUserLength == 0 {
		fmt.Printf("UserId:\t%v\tParents(Id):\t%v\n", elmentUserId, []int{0})
		return 0
	}

	for {
		for {
			rand.Seed(time.Now().UnixNano())
			if randomParentId = rand.Intn(elmentUserId-1) + 1; randomParentId < elmentUserId {
				break
			}
		}

		parents, stack := GetParentsStack(uint32(randomParentId))
		where["id"] = randomParentId
		parent = queries.SelectUser(where)

		if stack <= maxLengthTree {
			// TODO: Clear all fmt
			fmt.Printf("UserId:\t%v\tParents(Id):\t%v\n", elmentUserId, parents)
			break
		}
	}

	return parent.Id
}

func GetParentsStack(userId uint32) (parents []uint32, stack int) {
	var parent *dto.ResponseUser
	where := map[string]interface{}{}

	for {
		parents = append(parents, userId)
		where["id"] = userId
		parent = queries.SelectUser(where)
		stack++

		if parent.ParentId == 0 {
			break
		}

		userId = parent.ParentId
	}

	return parents, stack
}

func CheckChilndrensStackAndBallance(userId uint32, deep int) (_ [][]uint32, _ [][]uint32) {
	//First element will be userId -> length = deep+1
	childrensId := make([][]uint32, deep+1)
	childrensBalance := make([][]uint32, deep+1)
	i := 1

	childrensId[0] = append(childrensId[0], userId)
	childrensBalance[0] = append(childrensBalance[0], userId)

	//Checks all nesting levels and gives the previous level of childrenId
	for i = 1; i < deep+1; i++ {
		level, balance := CountingChildrensLevelAndBalance(childrensId[i-1])
		childrensId[i] = append(childrensId[i], level...)
		childrensBalance[i] = append(childrensBalance[i], balance...)
		if len(childrensId[i]) == 0 {
			break
		}
	}

	return childrensId[1:i], childrensBalance[1:i]
}

func CountingChildrensLevelAndBalance(parentsId []uint32) (level []uint32, balance []uint32) {
	for i := 0; i != len(parentsId); i++ {
		where := fmt.Sprint("parent_id = ", parentsId[i])
		childrens := queries.GetUsers(where)
		//For each element of parents[] get list of their childrens to -> childrens[]
		for j := 0; j != len(childrens); j++ {
			level = append(level, childrens[j].Id)
			balance = append(balance, childrens[j].Balance)
		}
	}

	return level, balance
}
