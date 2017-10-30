package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"os"
)

func initAllConnections() {
	messages := make(chan []byte)
	var urls [4]string

	urls[0] = "wss://api.gemini.com/v1/marketdata/btcusd"
	urls[1] = "ws://api.hitbtc.com:80"
	urls[2] = "wss://api.gemini.com/v1/marketdata/ethbtc"
	urls[3] = "wss://api.gemini.com/v1/marketdata/ethusd"

	for _, u := range urls {
		go func(u string) {
			ws, err := websocket.Dial(u, "", u)
			if err != nil {
				fmt.Println("FAILED CONNECTION")
				fmt.Printf("Dial failed: %s\n", err.Error())
				os.Exit(1)
			}

			defer ws.Close()

			incomingMessages := make(chan string)
			go readExchangeMessages(ws, incomingMessages)
			for {
				select {
				case message := <-incomingMessages:
					if u == urls[0] {
						fmt.Println(`Message Received FROM GEMINI BTCUSD:`, message)
					} else if u == urls[1] {
						fmt.Println(`Message Received FROM HITBTC:`, message)
					} else if u == urls[2] {
						fmt.Println(`Message Received FROM GEMINI ETHBTC:`, message)
					} else if u == urls[3] {
						fmt.Println(`Message Received FROM GEMINI ETHUSD:`, message)
					}
				}
			}
		}(u)
	}

	for m := range messages {
		fmt.Printf("%s\n", m)
	}
}

func readExchangeMessages(ws *websocket.Conn, incomingMessages chan string) {
	for {
		var message string
		// err := websocket.JSON.Receive(ws, &message)
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			fmt.Printf("Error Is: %s\n", err.Error())
			return
		}
		incomingMessages <- message
	}
}

func main() {
	initAllConnections()
}
