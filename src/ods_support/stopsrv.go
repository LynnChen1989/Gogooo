package main

import (
	"os"
)

const PAUSE = `{"sevendDayTime":"","sevendDayTimeName":"7天时间","updateInfo":"亲爱的用户，我行系统将于本日23点50分开始自动进行本日会计核算，持续时间约20分钟，给您带来的不便敬请谅解。","areaCodeTempName":"预加载地区 使  areacode 指向当前区域","areaCodeTemp":""}`

const RESTORE = `{"sevendDayTime":"","sevendDayTimeName":"7天时间","areaCodeTempName":"预加载地区 使  areacode 指向当前区域","areaCodeTemp":""}`

func PauseSrv() (status string) {
	client := RedisClient(os.Getenv("YHB_REDIS_HOST"), os.Getenv("YHB_REDIS_PASSWORD"), os.Getenv("YHB_REDIS_DB"))
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
	client := RedisClient(os.Getenv("YHB_REDIS_HOST"), os.Getenv("YHB_REDIS_PASSWORD"), os.Getenv("YHB_REDIS_DB"))
	err := client.Set("yhb_node_001_utilsRedisConfig", RESTORE, 0).Err()
	if err != nil {
		status = "failed"
		panic(err)
	}
	Info.Println("start srv success")
	status = "success"
	return
}
