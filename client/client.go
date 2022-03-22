package client

import (
	"encoding/json"
	"log"
	"main/config"
	"main/internal/trade"
	"net/url"
	"os"

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

func NewClient(scheme string, host string, path string) *Client {
	return &Client{
		Scheme: scheme, Host: host, Path: path,
	}
}
func (client Client) Connect(exchange *config.Exchange, responseChannel chan trade.Response, OSSignal chan os.Signal) {

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

	err = conn.WriteJSON(message)
	if err != nil {
		log.Println("Error during writing to websocket:", err)
		return
	}

	receiveHandler(conn, responseChannel)
	for sig := range OSSignal {
		log.Printf("Received %s signal. Closing all pending connections", sig)
		return
	}

}

func receiveHandler(connection *websocket.Conn, responseChannel chan trade.Response) {

	go func() {
		for {
			_, msg, err := connection.ReadMessage()
			if err != nil {
				log.Println("No messages to read; closing the response channel")
				close(responseChannel)
				return
			}
			resp := trade.Response{}
			err = json.Unmarshal(msg, &resp)
			if err != nil {
				log.Fatal("Fatal error occurred in unmarshalling")
				return
			}
			responseChannel <- resp
		}
	}()
}

func (client *Client) InitSignalHandler(OSSignal chan os.Signal) {

	sig := <-OSSignal
	log.Printf("Closing the channel with %s", sig)

	OSSignal <- sig
}
