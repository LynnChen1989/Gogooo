package main

import (
	"fmt"
	"os"
)

func SendNotify(level string, title string, content string) {

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
	message := fmt.Sprintf(`{"cid": %d, "title": "%s", "content": "%s"}`, messageCid, title, content)
	Info.Println(message)
	httpPost(messageUrl, message)
}
