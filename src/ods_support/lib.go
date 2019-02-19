package main

import (
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

func RedisClient(redisHost string, redisPassword string, redisDb string) (client *redis.Client) {
	//redisHost, redisPassword, redisDb := os.Getenv("YHB_REDIS_HOST"), os.Getenv("YHB_REDIS_PASSWORD"), os.Getenv("YHB_REDIS_DB")
	rDb, err := strconv.Atoi(redisDb)
	if err != nil {
		Error.Println("Covert variable redisDb to int error")
	}
	Info.Printf("redis connect info, redis host: %s, redis db: %d", redisHost, rDb)
	client = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword, // no password set
		DB:       rDb,           // use default DB
	})
	return
}

func NowFormatDate(format string) (rs string) {
	currentTime := time.Now()
	rs = currentTime.Format(format)
	return
}
