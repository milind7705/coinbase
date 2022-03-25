package trade

import (
	"log"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// Individual Trade from returned json.
type Trade struct {
	Price     decimal.Decimal `json:"price"`
	Size      decimal.Decimal `json:"size"`
	ProductId string          `json:"product_id"`
}

type Response struct {
	Type         string    `json:"type"`
	TradeID      int       `json:"trade_id"`
	MakerOrderID string    `json:"maker_order_id"`
	TakerOrderID string    `json:"taker_order_id"`
	Side         string    `json:"side"`
	Size         string    `json:"size"`
	Price        string    `json:"price"`
	ProductID    string    `json:"product_id"`
	Sequence     int64     `json:"sequence"`
	Time         time.Time `json:"time"`
}

// Queue to be used for first in first out and for sliding window
type Queue struct {
	MaxSize                int
	Lock                   *sync.Mutex
	Points                 []Trade
	SummationPriceQuantity map[string]decimal.Decimal
	SummationQuantity      map[string]decimal.Decimal
	VWAP                   map[string]decimal.Decimal
}

func NewQueue(MaxSize int) *Queue {
	return &Queue{
		MaxSize:                MaxSize,
		Lock:                   &sync.Mutex{},
		Points:                 []Trade{},
		SummationPriceQuantity: make(map[string]decimal.Decimal),
		SummationQuantity:      make(map[string]decimal.Decimal),
		VWAP:                   make(map[string]decimal.Decimal),
	}
}

func (q *Queue) Enqueue(element Trade) []Trade {
	q.Lock.Lock()
	defer q.Lock.Unlock()

	q.Points = append(q.Points, element)

	q.SummationPriceQuantity[element.ProductId] = q.SummationPriceQuantity[element.ProductId].Add(element.Price.Mul(element.Size))
	q.SummationQuantity[element.ProductId] = q.SummationQuantity[element.ProductId].Add(element.Size)
	q.VWAP[element.ProductId] = q.SummationPriceQuantity[element.ProductId].Div(q.SummationQuantity[element.ProductId])
	return q.Points
}

func (q *Queue) Dequeue() {
	q.Lock.Lock()
	defer q.Lock.Unlock()

	element := q.Points[0]
	q.Points = q.Points[1:]

	q.SummationPriceQuantity[element.ProductId] = q.SummationPriceQuantity[element.ProductId].Sub(element.Price.Mul(element.Size))
	q.SummationQuantity[element.ProductId] = q.SummationQuantity[element.ProductId].Sub(element.Size)
	if !q.SummationQuantity[element.ProductId].IsZero() {
		q.VWAP[element.ProductId] = q.SummationPriceQuantity[element.ProductId].Div(q.SummationQuantity[element.ProductId])
	}

}

func (q *Queue) Populate(responseChannel chan Response) {

	for point := range responseChannel {
		// first message send the empty price to the channel
		if point.Price == "" {
			continue
		}
		price, err := decimal.NewFromString(point.Price)
		if err != nil {
			log.Printf("error converting price %s: %v", point.Price, err)
			continue
		}

		size, err := decimal.NewFromString(point.Size)
		if err != nil {
			log.Printf("error converting size %s: %v", point.Size, err)
			continue
		}

		trade := Trade{Price: price, Size: size, ProductId: point.ProductID}

		if len(q.Points) == 0 {
			q.Points = append(q.Points, trade)
			q.SummationPriceQuantity[trade.ProductId] = trade.Price.Mul(trade.Size)
			q.SummationQuantity[trade.ProductId] = trade.Size
			q.VWAP[trade.ProductId] = trade.Price.Mul(trade.Size).Div(trade.Size)
		} else if len(q.Points) == q.MaxSize {
			q.Dequeue()
			q.Points = q.Enqueue(trade)
		} else {
			q.Points = q.Enqueue(trade)
		}
		log.Print(q.VWAP)
	}
}
