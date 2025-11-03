package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int8

const(
	Durable SimpleQueueType = 0
	Transient SimpleQueueType = 1
)


func main() {
	fmt.Println("Starting Peril client...")
	connstr:="amqp://guest:guest@localhost:5672/"
	// done := make(chan os.Signal, 1)
	conn, err:=amqp.Dial(connstr)
	if err== nil{
		fmt.Println("Connection Successfull")
	} else{
		fmt.Println("Connection Failed")
	}
	username,err:=gamelogic.ClientWelcome()
	if err!=nil{
		fmt.Printf("%s",err)
	}
	DeclareAndBind(conn,routing.ExchangePerilDirect,routing.PauseKey+"."+username,routing.PauseKey, 1)
	defer conn.Close()
	gamestate:=gamelogic.NewGameState(username)
	var words []string
	for{
		if words=gamelogic.GetInput();words==nil {
			continue
		}
		switch words[0] {
			case "spawn": {
				fmt.Println("Spawn...")
				gamestate.CommandSpawn(words)
				
			}
			case "move": {
				fmt.Println("Move...")
				gamestate.CommandMove(words)
			}
			case "status":{
				fmt.Println("Status...")
				gamestate.CommandStatus()
			}
			case "help":{
				fmt.Println("Help...")
				gamelogic.PrintClientHelp()
			}
			case "spam":{
				fmt.Println("Spamming not allowed yet!...")
			
			}
			case "quit":{
				gamelogic.PrintQuit()
				fmt.Println("Exiting Gracefully...")
				break
			}
			default: {
				fmt.Println("Unknown Command")
			}
		}
	}
		// signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
		// <-done
		// fmt.Println("Received signal, exiting gracefully...")
}
func DeclareAndBind(
		conn *amqp.Connection,
		exchange,
		queueName,
		key string,
		queueType SimpleQueueType, // an enum to represent "durable" or "transient"
		
	) (*amqp.Channel, amqp.Queue, error){
		ch,err := conn.Channel()
		if err!= nil{
			return ch,amqp.Queue{},err
		}
		var q amqp.Queue
		switch queueType {
			case 0:
				q,err= ch.QueueDeclare(queueName,true,false,false,false,nil)
			case 1:
				q,err= ch.QueueDeclare(queueName,false, true,true, false, nil)
			default:
				return ch, amqp.Queue{}, nil
		}
		if err!=nil{
			return ch, amqp.Queue{}, err
		}
		ch.QueueBind("amqpQueue",key,exchange,false,nil)
		
		return ch, q, nil
	}