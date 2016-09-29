package main

import (
	"fmt"
	"github.com/nlopes/slack"
)

func main() {
	// create new slack object & connect
	api := slack.New("TOKEN_REDACTED")
	api.SetDebug(true)
	rtm := api.NewRTM()
	// lets do it
	go rtm.ManageConnection()
	// loop() and listen for inbound messages
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			fmt.Printf("  %s\n", msg.Data)
		}
	}
}
