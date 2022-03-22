package client

import (
	"log"
	"main/config"
	"main/internal/trade"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// test config yaml for integration tests
const TestGoodConfig = "../config/test_exchange.yaml"

var upgrader = websocket.Upgrader{}

func tradeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()

	// send the sample trades on the socket
	response1 := `{"type": "last_match", "trade_id": 243091453, "maker_order_id": "fcb93eab-2d43-4856-b6d2-f20363bdac2d", "taker_order_id": "710c4c8e-3b0f-4bc6-93eb-79ced0fe6fa1", "side": "sell", "size": "0.01655371", "price": "2885.81", "product_id": "ETH-USD", "sequence": 27069295509, "time": "2022-03-20T14:19:42.998986Z"}`
	response2 := `{"type": "last_match", "trade_id": 243091454, "maker_order_id": "a8b80cff-0d4c-46b1-8bce-10ee39dee96e", "taker_order_id": "dd12e8b2-f80c-49ea-9f82-33dc4c8ca17d", "side": "sell", "size": "0.016553", "price": "2885.01", "product_id": "ETH-USD", "sequence": 27069295591, "time": "2022-03-20T14:19:42.998986Z"}`
	_ = c.WriteMessage(1, []byte(response1))
	_ = c.WriteMessage(1, []byte(response2))

}

func TestExample(t *testing.T) {
	// Integration test with creating a test server with the trade handler.
	s := httptest.NewServer(http.HandlerFunc(tradeHandler))
	defer s.Close()

	exchange, _ := config.NewExchange(TestGoodConfig)

	exchange.Host = strings.TrimPrefix(s.URL, "http://")
	log.Print(exchange.Host)

	client := NewClient(exchange.Scheme, exchange.Host, exchange.Path)

	queue := trade.NewQueue(exchange.Maxsize)

	// response channel to communicate between client and queue
	responseChannel := make(chan trade.Response)

	go queue.Populate(responseChannel)

	// channel for handling interrupts
	OSSignal := make(chan os.Signal, 1)

	signal.Notify(OSSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		time.Sleep(time.Second)
		OSSignal <- syscall.SIGINT
	}()

	client.Connect(exchange, responseChannel, OSSignal)
	// the value should match vwap of 2885.4100085783214339
	d, _ := decimal.NewFromString("2885.4100085783214339")
	assert.Equal(t, d.Cmp(queue.VWAP["ETH-USD"]), 0)

}
