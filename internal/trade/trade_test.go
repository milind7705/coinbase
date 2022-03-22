package trade

import (
	"main/config"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestQueue_Populate(t *testing.T) {
	t.Parallel()

	queue := NewQueue(config.Maxsize)

	responseChannel := make(chan Response)

	responses := []Response{
		{
			Type:         "match",
			TradeID:      300167031,
			MakerOrderID: "16de6b51-5001-4192-855e-71c2a6400616",
			TakerOrderID: "fdd909a3-ba72-4ce9-a862-599d88fa516f",
			Side:         "buy",
			Size:         "0.000601",
			Price:        "41630.9",
			ProductID:    "BTC-USD",
			Sequence:     35316892575,
		},
		{
			Type:         "match",
			TradeID:      300167032,
			MakerOrderID: "ed8e935a-7936-45d1-ab60-ed26f950dcac",
			TakerOrderID: "9b0a046f-211a-41c7-ba6b-97a98494c6f7",
			Side:         "buy",
			Size:         "0.00271786",
			Price:        "41631.45",
			ProductID:    "BTC-USD",
			Sequence:     35316892576,
		},
		{
			Type:         "match",
			TradeID:      300167033,
			MakerOrderID: "9a40e34f-1003-4d7a-b4d9-01e934015f64",
			TakerOrderID: "98018a07-1e4c-417f-a801-72aecbdb7202",
			Side:         "buy",
			Size:         "0.00005069",
			Price:        "41630.82",
			ProductID:    "BTC-USD",
			Sequence:     35316892577,
		},
	}

	// start the queue population on the response channel
	go queue.Populate(responseChannel)

	// send the sample responses on the channel
	for _, resp := range responses {
		responseChannel <- resp
	}

	close(responseChannel)
	time.Sleep(time.Millisecond)

	d, _ := decimal.NewFromString("41631.3424234096541081")

	assert.Equal(t, d.Cmp(queue.VWAP["BTC-USD"]), 0)
}
