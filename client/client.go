package client

import (
	"encoding/json"
	"log"
	"main/config"
	"main/internal/trade"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type       string           `json:"type"`
	ProductIds []string         `json:"product_ids"`
	Channels   []config.Channel `json:"channels"`
}

type Client struct {
	Scheme string
	Host   string
	Path   string
}

var done chan interface{}
var interrupt chan os.Signal

func NewClient(scheme string, host string, path string) *Client {
	return &Client{
		Scheme: scheme, Host: host, Path: path,
	}
}
func (client Client) Connect(exchange *config.Exchange) {

	done = make(chan interface{})

	interrupt = make(chan os.Signal)

	u := url.URL{Scheme: client.Scheme, Host: client.Host, Path: client.Path}

	message := Message{
		Type:       exchange.Message.Type,
		ProductIds: exchange.Message.ProductIds,
		Channels:   exchange.Message.Channels,
	}

	log.Printf("Connecting to %s", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()
	log.Print(message)

	err = conn.WriteJSON(message)
	if err != nil {
		log.Println("Error during writing to websocket:", err)
		return
	}

	go receiveHandler(conn)
	for {
		select {
		case <-interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			log.Println("Received SIGINT interrupt signal. Closing all pending connections")

			// Close our websocket connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}

			select {
			case <-done:
				log.Println("Receiver Channel Closed! Exiting....")
			case <-time.After(time.Duration(1) * time.Second):
				log.Println("Timeout in closing receiving channel. Exiting....")
			}
			return
		}
	}
}

func receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		resp := trade.Trade{}
		err = json.Unmarshal(msg, &resp)
		if err != nil {
			log.Fatal("Fatal")
		}
		log.Printf("Received: %v\n", resp)
	}
}
