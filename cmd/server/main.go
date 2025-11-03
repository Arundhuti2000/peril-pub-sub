package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	
	connstr:="amqp://guest:guest@localhost:5672/"
	// done := make(chan os.Signal, 1)
	conn, err:=amqp.Dial(connstr)
	if err== nil{
		fmt.Println("Connection Successfull")
	} else{
		fmt.Println("Connection Failed")
	}
	
	ch,err:=conn.Channel()
	if err== nil{
		fmt.Println("Channel Creation Successfull")
	} else{
		fmt.Println("Channel Creation Failed")
	}
	pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})
	defer conn.Close()

	pubsub.DeclareAndBind(conn,routing.ExchangePerilTopic,"game_logs","game_logs.*", 0)
	gamelogic.PrintServerHelp()
	var words []string
	for{
		if words=gamelogic.GetInput();words==nil {
			continue
		}
		switch words[0] {
			case "pause": {
				fmt.Println("Pause...")
				pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
					IsPaused: true,
				})
			}
			case "resume": {
				fmt.Println("Resume...")
				pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
					IsPaused: true,
				})
			}
			case "quit":{
				fmt.Println("Exiting Gracefully...")
				break
			}
			default: fmt.Println("Unknown Command")
		}
	}
	
	// signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	// <-done
	// fmt.Println("Received signal, exiting gracefully...")
	
}
