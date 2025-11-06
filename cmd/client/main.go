package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}

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
	// _,queue,err:=pubsub.DeclareAndBind(conn,routing.ExchangePerilDirect,routing.PauseKey+"."+username,routing.PauseKey, pubsub.Transient)
	defer conn.Close()
	// if err!= nil{
	// 	fmt.Println("Channel Declare and Bind failed")
	// } else{
	// 	fmt.Printf("Queue %v declared and bound!\n", queue.Name)
	// }
	gamestate:=gamelogic.NewGameState(username)
	handler:=handlerPause(gamestate)
	pubsub.SubscribeJSON(conn,routing.ExchangePerilDirect,routing.PauseKey+"."+username, routing.PauseKey,pubsub.Transient,handler)
	
	var words []string
	for{
		if words=gamelogic.GetInput();words==nil {
			fmt.Println(words)
			continue
		}
		fmt.Println(words)
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
