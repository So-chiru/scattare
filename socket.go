package main

import (
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var TWITCH_WSS_API = "irc-ws.chat.twitch.tv:443"

func connect(channel string, onmsg chan<- []byte) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: TWITCH_WSS_API, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if *DEBUG_MODE {
					log.Println("read:", err)
				}

				return
			}

			if strings.HasPrefix(string(message), "PING ") {
				var pong = strings.Replace(string(message), "PING", "PONG", 1)

				if *DEBUG_MODE {
					log.Println("received a ping message")
				}

				c.WriteMessage(websocket.TextMessage, []byte(pong))
				continue
			}

			onmsg <- message
		}
	}()

	c.WriteMessage(websocket.TextMessage, []byte("CAP REQ :twitch.tv/tags twitch.tv/commands"))
	c.WriteMessage(websocket.TextMessage, []byte("PASS SCHMOOPIIE"))

	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(99999-10000) + 10000
	c.WriteMessage(websocket.TextMessage, []byte("NICK justinfan"+strconv.Itoa(random)))

	c.WriteMessage(websocket.TextMessage, []byte("JOIN #"+channel))

	log.Printf("joining to #%s", channel)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			if *DEBUG_MODE {
				log.Println("interrupt")
			}

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

}
