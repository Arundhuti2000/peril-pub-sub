package main

import (
	"fmt"
	"os"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
)

func main() {
	fmt.Println("Starting Peril client...")
	connstr:="amqp://guest:guest@localhost:5672/"
	conn, err:=amqp.Dial(connstr)
	if err== nil{
		fmt.Println("Connection Successfull")
	} else{
		fmt.Println("Connection Failed")
	}
	gamelogic.ClientWelcome()
	func DeclareAndBind(
		conn *amqp.Connection,
		exchange,
		queueName,
		key string,
		queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	) (*amqp.Channel, amqp.Queue, error){
		ch := conn.Channel()
		if queueType == "durable"{
			q:= ch.QueueDeclare(
				durable: true
			)
		}
		
	}
}
