package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	<-done
	fmt.Println("Peril server closing..")
}
