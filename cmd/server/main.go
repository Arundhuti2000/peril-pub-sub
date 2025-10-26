package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connstr:="amqp://guest:guest@localhost:5672/"
	done := make(chan os.Signal, 1)
	conn, err:=amqp.Dial(connstr)
	
	if err== nil{
		fmt.Println("Connection Successfull")
	} else{
		fmt.Println("Connection Failed")
	}
	defer conn.Close()
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	fmt.Println("Received signal, exiting gracefully...")
	fmt.Println("Starting Peril server...")
}
