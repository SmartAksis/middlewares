package memory

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func InitRedis() {
	if rdb == nil {
		port, _ := strconv.ParseInt(os.Getenv("DB_REDIS_DATABASE"), 10, 64)
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", os.Getenv("DB_REDIS_HOST"), os.Getenv("DB_REDIS_PORT")),
			Password: os.Getenv("DB_REDIS_PASSWORD"), // no password set
			DB:       int(port),  // use default DB
		})
	}
}
//
//func SetParameterApplication(){
//
//	err := rdb.Set(ctx, "key", "value", 0).Err()
//	if err != nil {
//		panic(err)
//	}
//
//	val, err := rdb.Get(ctx, "key").Result()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("key", val)
//
//	val2, err := rdb.Get(ctx, "key2").Result()
//	if err == redis.Nil {
//		fmt.Println("key2 does not exist")
//	} else if err != nil {
//		panic(err)
//	} else {
//		fmt.Println("key2", val2)
//	}
//	// Output: key value
//	// key2 does not exist
//}
