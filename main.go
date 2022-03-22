package main

import (
	"log"
	"main/client"
	"main/config"
	"main/internal/trade"
	"os"
	"os/signal"
	"syscall"
)

const numberOfArgs = 2

var OSSignal chan os.Signal

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

	// response channel to communicate between client and queue
	responseChannel := make(chan trade.Response)

	go queue.Populate(responseChannel)

	// channel for handling interrupts
	OSSignal := make(chan os.Signal, 1)

	signal.Notify(OSSignal, syscall.SIGINT, syscall.SIGTERM)

	go client.InitSignalHandler(OSSignal)

	client.Connect(exchange, responseChannel, OSSignal)

}
