package main

import (
	"flag"
	"log"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var DEBUG_MODE = flag.Bool("debug", false, "enable debug mode")
var TRANSPORT_ENDPOINT = flag.String("e", "", "transport endpoint (accept https-http, leave empty to disable)")
var TRANSPORT_HEADERS = flag.String("h", "", "transport headers (serialized JSON)")
var OUPUT_FILE = flag.String("output", "data.json", "output file (.csv, .json)")
var COLLECT_INTERVAL = flag.Int("interval", 3000, "collect interval in miliseconds")

func main() {
	var channel = flag.String("channel", "", "channel to connect to")
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	if channel == nil || *channel == "" {
		panic("channel is required")
	}

	if OUPUT_FILE == nil || *OUPUT_FILE == "" {
		panic("output file is required")
	} else if !strings.HasSuffix(*OUPUT_FILE, ".csv") && !strings.HasSuffix(*OUPUT_FILE, ".json") {
		panic("output file must be .csv or .json")
	}

	if COLLECT_INTERVAL == nil || *COLLECT_INTERVAL < 10 || *COLLECT_INTERVAL > math.MaxInt32 {
		panic("interval must be between 10 and " + strconv.Itoa(math.MaxInt32))
	}

	ticker := time.NewTicker(time.Millisecond * time.Duration(*COLLECT_INTERVAL))

	go func() {
		for {
			select {
			case <-interrupt:
				log.Println("saving data before exiting...")
				save_worker()
				break
			case <-ticker.C:
				save_worker()
				break
			}
		}
	}()

	ch := make(chan []byte, 1)
	go func() {
		for {
			var raw = <-ch
			msg := parse(raw)

			if msg != nil {
				log.Println(msg.Channel, msg.Username, msg.Message)
			}

			stack_add(msg)
		}
	}()

	connect(*channel, ch)
}
