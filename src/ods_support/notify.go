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
	message := fmt.Sprintf(`{"cid": %s, "title": "%s", "content": "%s"}`, messageCid, title, content)
	httpPost(messageUrl, message)
}
