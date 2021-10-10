package main

import (
	"context"
	"net/http"
	"os"

	http_controllers "store/app/controllers"
	redis "store/app/redis"
	services "store/app/services"
	postgre "store/app/services/database/connect"
	fill "store/app/services/database/filldatabase"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	postgre.ConnectToDataBase()
	postgre.CreateDataBase()
	fill.FillDataBase(5, 7)

	services.GetWinners(10, 3, 15)
	go services.CronWinners(10, 3, 15)

	ctx := context.Background()
	redis.ConnectToRedis(ctx, os.Getenv("REDIS_port"))
	defer redis.Close()

	router := http_controllers.InitControllers()
	http.ListenAndServe(os.Getenv("SERVICE_PORT"), router)
}
