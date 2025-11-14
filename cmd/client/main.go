package main

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
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
	defer conn.Close()

	publishCh, err := conn.Channel()
	if err != nil {
		fmt.Println("Could not create publish channel")
		return
	}
	defer publishCh.Close()
	username,err:=gamelogic.ClientWelcome()
	if err!=nil{
		fmt.Printf("%s",err)
	}
	// _,queue,err:=pubsub.DeclareAndBind(conn,routing.ExchangePerilDirect,routing.PauseKey+"."+username,routing.PauseKey, pubsub.Transient)
	

	
	// if err!= nil{
	// 	fmt.Println("Channel Declare and Bind failed")
	// } else{
	// 	fmt.Printf("Queue %v declared and bound!\n", queue.Name)
	// }
	gamestate:=gamelogic.NewGameState(username)
	handler:=handlerPause(gamestate)
	acktype,err:=pubsub.SubscribeJSON(conn,routing.ExchangePerilDirect,routing.PauseKey+"."+username, routing.PauseKey,pubsub.Transient,handler)
	if err!=nil{
		fmt.Printf("%s",err)
	}
	
	handlerMove:=handlerMove(gamestate)
	pubsub.SubscribeJSON(conn,routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+username,routing.ArmyMovesPrefix+".*", pubsub.Transient,handlerMove)
	var words []string
	for{
		if words=gamelogic.GetInput();len(words)==0 {
			fmt.Println(words)
			continue
		}
		fmt.Println(words)
		switch words[0] {
			case "spawn": {
				fmt.Println("Spawn...")
				err:=gamestate.CommandSpawn(words)
				if err != nil {
				fmt.Println(err)
				continue
			}
				
			}
			case "move": {
				fmt.Println("Move...")
				move,err:=gamestate.CommandMove(words)
				if err != nil {
					fmt.Println(err)
					continue
				}
				pubsub.PublishJSON(publishCh,routing.ExchangePerilTopic,routing.ArmyMovesPrefix+"."+username,move)
				err = pubsub.PublishJSON(publishCh, routing.ExchangePerilTopic, routing.GameLogSlug+"."+username, routing.GameLog{
					CurrentTime: time.Now(),
					Message:     fmt.Sprintf("%s moved units to %s", username, move.ToLocation),
					Username:    username,
				})
				if err != nil {
					fmt.Println("Error publishing game log:", err)
				}
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
