package redis

import (
	"context"
	"log"
	"time"

	voc "store/app/vocabulary"

	"github.com/go-redis/redis/v8"
)

var Session *RedisSession

var client *redis.Client

type RedisSession struct {
	Client *redis.Client
	Ctx    context.Context
}

func ConnectToRedis(ctx context.Context, addr string) {
	var client *redis.Client

	for i := 0; i < 3; i++ {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
		})

		if client != nil {
			break
		}

		log.Println(voc.ERROR_CONNECT_REDIS_RECONNECT)
		time.Sleep(time.Duration(i*5) * time.Second)
	}

	if client == nil {
		log.Println(voc.ERROR_CONNECT_REDIS)
		panic(voc.ERROR_CONNECT_REDIS)
	}

	// iter := client.Scan(ctx, 0, "*", 0).Iterator()
	// for iter.Next(ctx) {
	// 	val := iter.Val()
	// 	err := client.Del(ctx, val).Err()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// if err := iter.Err(); err != nil {
	// 	panic(err)
	// }

	Session = &RedisSession{
		Ctx:    ctx,
		Client: client,
	}
}

func Close() {
	Session.Client.Close()
}
