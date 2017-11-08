package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"os"
)

func initAllConnections() {
	messages := make(chan []byte)
	var urls [4]string

	urls[0] = "ws://api.hitbtc.com:80"
	urls[1] = "wss://api.gemini.com/v1/marketdata/btcusd"
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
			go readExchangeMessages(ws, incomingMessages, u)

			var countNum int64
			countNum = 0

			for {
				select {
				case message := <-incomingMessages:
					if u == urls[0] {
						countNum = countNum + 1
						fmt.Printf("\n\n\n\nMessage Received FROM HITBTC: %s", message)
						fmt.Printf("\n\n NUMBER OF Polls ARE: %d \n\n\n", countNum)
					} else if u == urls[1] {
						countNum = countNum + 1
						fmt.Println(`Message Received FROM GEMINI BTCUSD:`, message)
						fmt.Printf("\n\n NUMBER OF Polls ARE: %d \n\n\n", countNum)
					} else if u == urls[2] {
						countNum = countNum + 1
						fmt.Println(`Message Received FROM GEMINI ETHBTC:`, message)
						fmt.Printf("\n\n NUMBER OF Polls ARE: %d \n\n\n", countNum)
					} else if u == urls[3] {
						countNum = countNum + 1
						fmt.Println(`Message Received FROM GEMINI ETHUSD:`, message)
						fmt.Printf("\n\n NUMBER OF Polls ARE: %d \n\n\n", countNum)
					}
				}
			}
		}(u)
	}

	for m := range messages {
		fmt.Printf("%s\n\n\n\n\n", m)
	}
}

func readExchangeMessages(ws *websocket.Conn, incomingMessages chan string, u string) {
	for {
		hitbtcResponse := new(struct {
			MarketDataIncrementalRefresh struct {
				SeqNo          int    `json:"seqNo"`
				Symbol         string `json:"symbol"`
				ExchangeStatus string `json:"exchangeStatus"`
				Ask            []struct {
					Price string `json:"price"`
					Size  int    `json:"size"`
				} `json:"ask"`
				Bid []struct {
					Price string `json:"price"`
					Size  int    `json:"size"`
				} `json:"bid"`
				Trade []struct {
					Price     string `json:"price"`
					Size      int    `json:"size"`
					TradeID   int    `json:"tradeId"`
					Timestamp int64  `json:"timestamp"`
					Side      string `json:"side"`
				} `json:"trade"`
				Timestamp int64 `json:"timestamp"`
			} `json:"MarketDataIncrementalRefresh"`
		})

		// geminiResponse := new(struct {
		// 	Type           string `json:"type"`
		// 	EventID        int64  `json:"eventId"`
		// 	Timestamp      int    `json:"timestamp"`
		// 	Timestampms    int64  `json:"timestampms"`
		// 	SocketSequence int    `json:"socket_sequence"`
		// 	Events         []struct {
		// 		Type      string `json:"type"`
		// 		Side      string `json:"side"`
		// 		Price     string `json:"price"`
		// 		Remaining string `json:"remaining"`
		// 		Delta     string `json:"delta"`
		// 		Reason    string `json:"reason"`
		// 	} `json:"events"`
		// })

		var message string
		// type jsonMessage struct {
		// 	Data string `json:"data"`
		// }

		// var jj jsonMessage
		// err := websocket.JSON.Receive(ws, &jj)

		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			fmt.Printf("Error Is: %s\n", err.Error(), err)
			return
		}

		if u == "ws://api.hitbtc.com:80" {
			if err := json.NewDecoder(ws).Decode(hitbtcResponse); err != nil {
				log.Fatal(err)
			}

			hitbtcMessage := map[string]interface{}{"symbol": hitbtcResponse.MarketDataIncrementalRefresh.Symbol, "exchangeStatus": hitbtcResponse.MarketDataIncrementalRefresh.ExchangeStatus, "ask": hitbtcResponse.MarketDataIncrementalRefresh.Ask, "bid": hitbtcResponse.MarketDataIncrementalRefresh.Bid, "trade": hitbtcResponse.MarketDataIncrementalRefresh.Trade, "timestamp": hitbtcResponse.MarketDataIncrementalRefresh.Timestamp}

			fmt.Printf("NEW INTERFACE IS: %s", hitbtcMessage)
		}
		//  else {
		// 	if err := json.NewDecoder(ws).Decode(geminiResponse); err != nil {
		// 		log.Fatal(err)
		// 	}

		// 	geminiMessage := map[string]interface{}{"type": geminiResponse.Events.Type, "side": geminiResponse.Events.Side, "price": geminiResponse.Events.Price, "timestamp": geminiResponse.Timestamp}

		// 	fmt.Printf("NEW INTERFACE IS: %s", geminiMessage)
		// }

		incomingMessages <- message
	}
}

func main() {
	initAllConnections()
}
