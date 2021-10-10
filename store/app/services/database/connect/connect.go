package connect

import (
	"fmt"
	"os"

	model "store/app/services/database/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ConncetDataBase *gorm.DB

func ConnectToDataBase() {
	dsn := fmt.Sprintf("host=%s user=%s database=%s password=%s port=%s",
		os.Getenv("POSTGRESQL_host"),
		os.Getenv("POSTGRESQL_user"),
		os.Getenv("POSTGRESQL_database"),
		os.Getenv("POSTGRESQL_password"),
		os.Getenv("POSTGRESQL_port"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: true})

	if err != nil {
		panic("failed to connect database")
	}

	ConncetDataBase = db
	db.Exec("CREATE TYPE StatusEnum AS ENUM('1', '2');")
}

func CreateDataBase() {
	ConncetDataBase.AutoMigrate(&model.User{})
	ConncetDataBase.AutoMigrate(&model.Transaction{})
	ConncetDataBase.AutoMigrate(&model.Category{})
	ConncetDataBase.AutoMigrate(&model.Product{})
}
