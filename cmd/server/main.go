package main

import (
	"bytes"
	"encoding/gob"
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

	// Subscribe to all game log messages and write them to disk
	err = pubsub.SubscribeGob(
		conn,
		routing.ExchangePerilTopic,
		"game_logs",
		routing.GameLogSlug+".*",
		pubsub.Durable,
		func(gl routing.GameLog) pubsub.Acktype {
			defer fmt.Print("> ")
			if err := gamelogic.WriteLog(gl); err != nil {
				fmt.Println("failed to write log:", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		},
		func(data []byte) (routing.GameLog, error) {
			var gl routing.GameLog
			dec := gob.NewDecoder(bytes.NewReader(data))
			err := dec.Decode(&gl)
			return gl, err
		},
	)
	if err != nil {
		fmt.Println("Failed to subscribe to game logs:", err)
	} else {
		fmt.Println("Subscribed to game logs")
	}
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
					IsPaused: false,
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
