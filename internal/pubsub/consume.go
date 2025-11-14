package pubsub

import (
	"encoding/json"
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
)
type SimpleQueueType int8
type acktype int8


const(
	Durable SimpleQueueType = 0
	Transient SimpleQueueType = 1
)

const (
	Ack acktype = 1
	NackRequeue acktype = 2
	NackDiscar acktype = 0
)

// func UnMarshal[T any](chDeli <-chan amqp.Delivery) amqp.Delivery{
// 	for val := range chDeli{
// 		var result []T
// 		json.Unmarshal(val.Body,result)
// 	}
// }
func SubscribeJSON[T any](
    conn *amqp.Connection,
    exchange,
    queueName,
    key string,
    queueType SimpleQueueType, // an enum to represent "durable" or "transient"
    handler func(T) acktype,
) error{
	ch, queue,err:=DeclareAndBind(conn,exchange,queueName,key,queueType)
	if err!=nil{
		return err
	}
	// 
	// channel:= make(chan amqp.Delivery)
	// channel.Consume()
	// channel,err:=conn.Channel()
	
	
	chDeli,err:=ch.Consume(queue.Name,"",false,false,false,false,nil)
	if err!=nil{
		return  fmt.Errorf("could not consume messages: %v", err)
	}
	
	
	// go UnMarshal[T](chDeli)
	go func() {
		defer ch.Close()
		for val := range chDeli {
			var result T
			json.Unmarshal(val.Body, &result)  
			acktype:=handler(result)
			switch acktype{
			case pubsub.Ack:
				val.Ack(false)
			case pubsub.NackRequeue:
				val.Nack(false,true)
			case pubsub.NackDiscar:
				val.Nack(false,false)
			}
			
		}
	}()
	return nil
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
			return ch,amqp.Queue{},fmt.Errorf("could not create channel: %v", err)
		}
		
		var q amqp.Queue
		switch queueType {
			case 0:
				q,err= ch.QueueDeclare(queueName,true,false,false,false,nil)
			case 1:
				q,err= ch.QueueDeclare(queueName,false, true,true, false, nil)
			default:
				return ch, amqp.Queue{}, fmt.Errorf("could not find queue type: %v", err)
		}
		if err!=nil{
			return ch, amqp.Queue{}, fmt.Errorf("could not declare: %v", err)
		}
		err=ch.QueueBind(queueName,key,exchange,false,nil)
		if err!=nil{
			return nil, amqp.Queue{},fmt.Errorf("could not bind queue: %v", err)
		}
		
		return ch, q, nil
	}