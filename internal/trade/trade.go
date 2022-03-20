package trade

import (
	"fmt"
	"sync"

	"github.com/shopspring/decimal"
)

const MaxSize = 200

// Individual Trade from returned json.
type Trade struct {
	Price     decimal.Decimal `json:"price"`
	Size      decimal.Decimal `json:"size"`
	ProductId string          `json:"product_id"`
}

// Queue to be used for first in first out and for sliding window
type Queue struct {
	MaxSize int

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

func Enqueue(p []Trade, element Trade) []Trade {
	p = append(p, element) // Simply append to enqueue.
	fmt.Println("Enqueued:", element)
	return p
}

func Dequeue(p []Trade) []Trade {
	element := p[0] // The first element is the one to be dequeued.
	fmt.Println("Dequeued:", element)
	return p[1:] // Slice off the element once it is dequeued.
}

func (q Queue) Populate() {

}
