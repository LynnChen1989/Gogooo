package main

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

var conn *amqp.Connection
var channel *amqp.Channel
var count = 0

const queueName = "snake.queue.test"
const exchange = "snakehost.exchange.direct"
const uri = "amqp://snake:snake@127.0.0.1:5672/snakehost"

func main() {
	go func() {
		for {
			push()
			time.Sleep(1 * time.Second)
		}
	}()
	receive()
	fmt.Println("end")
	mclose()
}
func failOnErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
		panic(fmt.Sprintf("%s:%s", msg, err))
	}
}
func mqConnect() {
	var err error
	conn, err = amqp.Dial(uri)
	failOnErr(err, "failed to connect tp rabbitmq")

	channel, err = conn.Channel()
	failOnErr(err, "failed to open a channel")
}

func mclose() {
	channel.Close()
	conn.Close()
}

//连接rabbitmq server
func push() {

	if channel == nil {
		mqConnect()
	}
	msgContent := "hello world!"

	channel.Publish(exchange, queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msgContent),
	})
}

func receive() {
	if channel == nil {
		mqConnect()
	}
	channel.QueueBind(queueName, queueName, exchange, false, nil)
	msgs, err := channel.Consume(queueName, "", true, false, false, false, nil)
	failOnErr(err, "")

	forever := make(chan bool)

	go func() {
		//fmt.Println(*msgs)
		for d := range msgs {
			s := BytesToString(&(d.Body))
			count++
			fmt.Printf("receve msg is :%s -- %d\n", *s, count)
		}
	}()

	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-forever
}

func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}
