package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"os"
)

var conn *amqp.Connection
var channel *amqp.Channel

const exchange = "ops.event.exchange"
const odsQueue = "ooops.batch.finish.queue"
const pushKey = "loan.batch.finish"
const receiveKey = "ops.listener.replicate"

func mqConnect() {
	var err error
	uri := os.Getenv("MQ_URI")
	conn, err = amqp.Dial(uri)
	if err != nil {
		Error.Println("failed to connect tp rabbitmq")
	}

	channel, err = conn.Channel()
	if err != nil {
		Error.Println("failed to open a channel")
	}
}

func mclose() {
	channel.Close()
	conn.Close()
}

func PushCutBatchMsg() {
	type CutStatus struct {
		Dt     string `json:"dt"`
		Status string `json:"status"`
	}
	if channel == nil {
		mqConnect()
	}
	today := NowFormatDate("2006-01-02")
	cs := &CutStatus{
		Dt:     today,
		Status: "finish",
	}
	message, _ := json.Marshal(cs)
	//msgContent := "hello world!"
	Info.Printf("pushing message: [%s]", message)
	channel.Publish(exchange, pushKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        message,
	})
}

func ReceiveMsg() {
	if channel == nil {
		mqConnect()
	}
	Info.Printf("declare queue name, ods: [%s], ", odsQueue)
	channel.QueueDeclare(odsQueue, true, true, false, true, nil)
	Info.Printf("bind queue [%s] to exchange [%s]", odsQueue, exchange)
	channel.QueueBind(odsQueue, receiveKey, exchange, false, nil)
	Info.Printf("Queue [%s] Waiting for MQ messages", odsQueue)
	messages, err := channel.Consume(odsQueue, "", true, false, false, false, nil)
	if err != nil {
		Error.Println("consume data error:", err)
	}
	go func() {
		var odsFinish map[string]interface{}
		for d := range messages {
			if err := json.Unmarshal(d.Body, &odsFinish); err == nil {
				Info.Printf("receive mq message, system:[%s], status:[%s], cutdate:[%s]",
					odsFinish["system"].(string), odsFinish["status"].(string), odsFinish["cutDate"].(string))
				//WriteOdsStatus(odsFinish["system"].(string), odsFinish["status"].(string))

				// 写入各个下游系统的状态到redis
				key := NowFormatDate("20060102") + odsFinish["system"].(string)

				// 判断获取到当日的日期
				if NowFormatDate("2006-01-02") == odsFinish["cutDate"].(string) {
					client := RedisClient(
						os.Getenv("CACHE_REDIS_HOST"),
						os.Getenv("CACHE_REDIS_PASSWORD"),
						os.Getenv("CACHE_REDIS_DB"))
					client.Set(key, odsFinish["status"].(string), 0)
				} else {
					Error.Println("fuck")
				}
			}
		}
	}()

}

func GetTodayOdsAllStatusFormRedis() (status bool) {
	status = false
	var success int
	var failure int
	systems := [3]string{"ods", "bigdata", "python"}
	for _, sys := range systems {
		key := NowFormatDate("20060102") + sys
		client := RedisClient(
			os.Getenv("CACHE_REDIS_HOST"),
			os.Getenv("CACHE_REDIS_PASSWORD"),
			os.Getenv("CACHE_REDIS_DB"))
		val := client.Get(key)
		if val.Val() == "failure" {
			failure += 1
		} else if val.Val() == "success" {
			success += 1
		}
	}
	if success == 3 && failure == 0 {
		status = true
	}
	return
}

//func BytesToString(b *[]byte) *string {
//	s := bytes.NewBuffer(*b)
//	r := s.String()
//	return &r
//}

//func WriteOdsStatus(system string, status string) {
//	judgeDbInfo := os.Getenv("JUDGE_ODS_DB_INFO")
//	dbw := DbWorker{Dsn: judgeDbInfo}
//	db, err := sql.Open("mysql", dbw.Dsn)
//	err = db.Ping()
//	if err != nil {
//		Error.Printf("Failed to ping mysql: %s", err)
//	}
//	if err != nil {
//		panic(err)
//		return
//	}
//	smt := fmt.Sprintf("insert into tb_restore_slave(system,status) values ('%s', '%s');", system, status)
//	if _, err := db.Query(smt); err == nil {
//		Info.Printf("write system [%s] status [%s] to db success", system, status)
//	}
//	db.Close()
//}

//
//func GetTodayOdsAllStatus() (status bool) {
//	var cnt int
//	currentTime := time.Now()
//	judgeDbInfo := os.Getenv("JUDGE_ODS_DB_INFO")
//	dbw := DbWorker{Dsn: judgeDbInfo}
//	db, err := sql.Open("mysql", dbw.Dsn)
//	err = db.Ping()
//	if err != nil {
//		Error.Printf("Failed to ping mysql: %s", err)
//	}
//	if err != nil {
//		panic(err)
//		return
//	}
//	smt01 := fmt.Sprintf("select count(1) as cnt from (select distinct * from tb_restore_slave where dt = '%s' group by system) a;",
//		currentTime.Format("2006-01-02"))
//	rows := db.QueryRow(smt01)
//	db.Close()
//	rows.Scan(&cnt)
//	if cnt == 3 {
//		status = true
//		Info.Printf("all system finish their ods job, cnt: %d,success", cnt)
//	} else {
//		status = false
//		Warning.Printf("some system have not finish their ods job, cnt: %d, failure", cnt)
//	}
//	return
//}
