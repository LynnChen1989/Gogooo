package main

import (
	"github.com/go-redis/redis"
	//"github.com/robfig/cron"
	"os"
	"strconv"
)

const PAUSE = `{"sevendDayTime":"","sevendDayTimeName":"7天时间","updateInfo":"亲爱的用户，我行将于本日00点00分进行会计核算，持续时间约20分钟，给您带来的不便敬请谅解。","areaCodeTempName":"预加载地区 使  areacode 指向当前区域","areaCodeTemp":""}`

const RESTORE = `{"sevendDayTime":"","sevendDayTimeName":"7天时间","areaCodeTempName":"预加载地区 使  areacode 指向当前区域","areaCodeTemp":""}`

func RedisClient() (client *redis.Client) {
	redisHost, redisPassword, redisDb := os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_DB")
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

func PauseSrv() (status string) {
	client := RedisClient()
	err := client.Set("yhb_node_001_utilsRedisConfig", PAUSE, 0).Err()
	if err != nil {
		status = "failed"
		panic(err)
	}
	Info.Println("stop srv success")
	status = "success"
	return
}

func RestoreSrv() (status string) {
	client := RedisClient()
	err := client.Set("yhb_node_001_utilsRedisConfig", RESTORE, 0).Err()
	if err != nil {
		status = "failed"
		panic(err)
	}
	Info.Println("start srv success")
	status = "success"
	return
}
