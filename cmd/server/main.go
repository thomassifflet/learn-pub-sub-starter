package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	sigs := make(chan os.Signal, 1)
	connectionString := "amqp://guest:guest@localhost:5672/"
	connectionDial, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Errorf("error: %v", err)
		return
	}

	qpChannel, err := connectionDial.Channel()
	if err != nil {
		fmt.Errorf("error creating channel: %v", err)
	}

	defer connectionDial.Close()
	fmt.Println("Connection to the RabbitMQ server successful.")

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("Starting Peril server...")
	gamelogic.PrintServerHelp()
	for {
		input := gamelogic.GetInput()
		if input == nil {
			continue
		}
		if input[0] == "pause" {
			fmt.Println("sending a pause message..")
			err = pubsub.PublishJSON(
				qpChannel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: true,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
				return
			}
		}
		if input[0] == "resume" {
			fmt.Println("sending a resume message..")
			err = pubsub.PublishJSON(
				qpChannel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
				return
			}
		}
		if input[0] == "quit" {
			fmt.Println("Exiting..")
			break
		}
		fmt.Println("command not understood")
	}
	<-done
	fmt.Println("Peril server closing..")
}
