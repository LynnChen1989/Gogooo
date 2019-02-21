package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

func RedisClient() (client *redis.Client) {
	redisHost, redisPassword, redisDb := "127.0.0.1:6379", "fuckyou", "15"
	rDb, err := strconv.Atoi(redisDb)
	if err != nil {
		fmt.Println("Covert variable redisDb to int error")
	}
	//fmt.Printf("redis connect info, redis host: %s, redis db: %d", redisHost, rDb)
	client = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword, // no password set
		DB:       rDb,           // use default DB
	})
	return
}
func main() {
	//currentTime := time.Now()
	//key := currentTime.Format("20060102") + "cutend"
	////fmt.Println(key)
	//client := RedisClient()
	//client.Set(key, "ok", 10000000)
	fmt.Println("a")
	time.Sleep(time.Second * 10)
	fmt.Println("b")
}
