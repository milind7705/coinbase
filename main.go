package main

import (
	"log"
	"main/client"
	"main/config"
	"main/internal/trade"
	"os"
)

const numberOfArgs = 2

func main() {

	exchange, err := config.DefaultExchange()
	if err != nil {
		panic("Unable to create exchange, check the default configs.")
	}

	if len(os.Args) != numberOfArgs {
		log.Printf("Config file missing, using default args to initialize.")
	} else {
		exchange, err = config.NewExchange(os.Args[1])

		if err != nil {
			panic("Unable to create exchange, check the config yaml.")
		}
	}

	client := client.NewClient(exchange.Scheme, exchange.Host, exchange.Path)

	queue := trade.NewQueue(exchange.Maxsize)

	responseChannel := make(chan trade.Response)

	go queue.Populate(responseChannel)

	client.Connect(exchange, responseChannel)

	go client.InitSignalHandler()
}
