package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"time"
)

var conn *amqp.Connection
var channel *amqp.Channel

const exchange = "ops.event.exchange"
const odsQueue = "ods.batch.finish.queue"
const bigDataQueue = "bigdata.batch.finish.queue"
const pythonQueue = "python.batch.finish.queue"

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

func pushMsg() {
	if channel == nil {
		mqConnect()
	}
	msgContent := "hello world!"

	channel.Publish(exchange, "loan.batch.finish", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msgContent),
	})
}

func receiveMsg() {
	if channel == nil {
		mqConnect()
	}
	Info.Printf("declare queue name, ods: [%s], bigdata: [%s], python: [%s]", odsQueue, bigDataQueue, pythonQueue)
	queues := [3]string{odsQueue, bigDataQueue, pythonQueue}
	for _, q := range queues {
		channel.QueueDeclare(q, true, true, false, true, nil)
		Info.Printf("bind queue [%s] to exchange [%s]", q, exchange)
		channel.QueueBind(q, q, exchange, false, nil)
		Info.Printf(" Queue [%s] Waiting for MQ messages", q)
		messages, err := channel.Consume(q, "", true, false, false, false, nil)
		if err != nil {
			Error.Println("consume data error:", err)
		}
		// 消息接受后处理
		go func() {
			var odsFinish map[string]interface{}
			for d := range messages {
				if err := json.Unmarshal(d.Body, &odsFinish); err == nil {
					fmt.Println(odsFinish["status"], odsFinish["system"])
					writeOdsStatus(odsFinish["system"].(string), odsFinish["status"].(string))

					if getTodayOdsAllStatus() {
						Info.Println("All ods job done, restore slave")
						//BatchHandleDbSlave("start")
						//TODO:上线开启
					}
				}
			}
		}()
	}

}

func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}

func writeOdsStatus(system string, status string) {
	judgeDbInfo := os.Getenv("JUDGE_ODS_DB_INFO")
	dbw := DbWorker{Dsn: judgeDbInfo}
	db, err := sql.Open("mysql", dbw.Dsn)
	err = db.Ping()
	if err != nil {
		Error.Printf("Failed to ping mysql: %s", err)
	}
	if err != nil {
		panic(err)
		return
	}
	smt := fmt.Sprintf("insert into tb_restore_slave(system,status) values ('%s', '%s');", system, status)
	if _, err := db.Query(smt); err == nil {
		Info.Printf("write system [%s] status [%s] to db success", system, status)
	}
	db.Close()
}

func getTodayOdsAllStatus() (status bool) {
	var cnt int
	currentTime := time.Now()
	judgeDbInfo := os.Getenv("JUDGE_ODS_DB_INFO")
	dbw := DbWorker{Dsn: judgeDbInfo}
	db, err := sql.Open("mysql", dbw.Dsn)
	err = db.Ping()
	if err != nil {
		Error.Printf("Failed to ping mysql: %s", err)
	}
	if err != nil {
		panic(err)
		return
	}
	smt01 := fmt.Sprintf("select count(1) as cnt from (select distinct * from tb_restore_slave where dt = '%s' group by system) a;",
		currentTime.Format("2006-01-02"))
	rows := db.QueryRow(smt01)
	db.Close()
	rows.Scan(&cnt)
	if cnt == 3 {
		status = true
		Info.Printf("all system finish their ods job, cnt: %d,success", cnt)
	} else {
		status = false
		Warning.Printf("some system have not finish their ods job, cnt: %d, failure", cnt)
	}
	return
}
