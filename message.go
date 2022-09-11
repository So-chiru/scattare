package main

import (
	"log"
	"strings"
)

type Message struct {
	Message  string            `json:"message"`
	Channel  string            `json:"channel"`
	Username string            `json:"username"`
	Tags     map[string]string `json:"tags"`
}

// parse message from twitch irc ws
func parse(raw []byte) *Message {
	if len(raw) == 0 {
		return nil
	}

	var msg = string(raw)

	if *DEBUG_MODE {
		log.Println("recv", msg)
	}

	var data = strings.Split(msg, " ")

	if len(data) < 3 {
		log.Printf("invalid message: %s", msg)

		return nil
	}

	if data[1] == "JOIN" {
		log.Printf("joined to channel %s", data[2])
		return nil
	}

	if !strings.HasPrefix(msg, "@badge-info") {
		return nil
	}

	var rawTags = strings.Split(data[0], ";")

	var tags = make(map[string]string)
	for _, tag := range rawTags {
		var tagData = strings.Split(tag, "=")
		tags[tagData[0]] = tagData[1]
	}

	var user = strings.Split(data[1], "!")
	// var channel = strings.Split(data[2], "#")

	var message = strings.Join(data[4:], " ")

	if message == "" {
		return nil
	}

	message = strings.ReplaceAll(message[1:], "\r\n", "")

	var username = strings.Split(user[0], ":")[1]

	return &Message{
		Channel:  data[3],
		Message:  message,
		Username: username,
		Tags:     tags,
	}
}
