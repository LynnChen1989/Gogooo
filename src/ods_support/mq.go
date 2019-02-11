package main

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel

//var count = 0
const queueName = "snake.queue.test"
const exchange = "ods.finish.status"
const uri = "amqp://snake:snake@127.0.0.1:5672/snakehost"

func mqConnect() {
	var err error
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

	channel.Publish(exchange, queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msgContent),
	})
}

func receiveMsg() {
	if channel == nil {
		mqConnect()
	}
	//channel.ExchangeDeclare(exchange,"topic", true, false, false,true,nil)
	Info.Printf("declare queue name: [%s]", queueName)
	channel.QueueDeclare(queueName, true, true, false, true, nil)
	Info.Printf("bind queue [%s] to exchange [%s]", queueName, exchange)
	channel.QueueBind(queueName, queueName, exchange, false, nil)
	messages, err := channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		Error.Println("consume data error:", err)
	}

	forever := make(chan bool)

	go func() {
		//fmt.Println(*msgs)
		for d := range messages {
			s := BytesToString(&(d.Body))
			fmt.Printf("receve msg is :%s\n", *s)
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
