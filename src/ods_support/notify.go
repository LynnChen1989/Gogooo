package main

import (
	"fmt"
	"os"
)

type NotifyRequest struct {
	cid     string
	title   string
	content string
}

func SendNotify(level string, title string, content string) {
	if os.Getenv("MESSAGE_URL") == "" {
		Error.Println("environment variable MESSAGE_URL is needed")
	}

	var messageCid int
	// TODO: 这里后面根据级别调整对应的告警CID
	if level == "general" {
		messageCid = 10004
	} else if level == "average" {
		messageCid = 10004
	} else if level == "fatal" {
		messageCid = 10004
	} else {
		Error.Println("Notify level error, please assign level of [general, average, fatal]")
		return
	}

	messageUrl := os.Getenv("MESSAGE_URL")
	Info.Printf("MESSAGE_URL is: %s", messageUrl)
	message := fmt.Sprintf(`{"cid": %d, "title": "%s", "content": "%s"}`, messageCid, title, content)
	Info.Println(message)
	httpPost(messageUrl, message)
}
