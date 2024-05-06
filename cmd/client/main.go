package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game client connected to RabbitMQ!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gameState := gamelogic.NewGameState(username)
	gamelogic.PrintClientHelp()
	for {
		input := gamelogic.GetInput()
		if input == nil {
			continue
		}
		if input[0] == "spawn" {
			err := gameState.CommandSpawn(input)
			if err != nil {
				log.Printf("couldnt execute command: %v", err)
				return
			}
		}
		if input[0] == "move" {
			_, err := gameState.CommandMove(input)
			if err != nil {
				log.Printf("couldnt execute command, error : %v", err)
				return
			}
			fmt.Println("Move successful")
		}
		if input[0] == "status" {
			gameState.CommandStatus()
		}
		if input[0] == "help" {
			gamelogic.PrintClientHelp()
		}
		if input[0] == "quit" {
			fmt.Println("exiting..")
			break
		}
		if input[0] == "spam" {
			fmt.Println("Spamming not allowed yet!")
		}
		if input[0] != "quit" && input[0] != "help" && input[0] != "move" && input[0] != "spawn" && input[0] != "status" && input[0] != "spam" {
			fmt.Println("Unknown command")
			continue
		}

	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
