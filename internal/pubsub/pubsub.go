package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
)
type SimpleQueueType int8

const(
	Durable SimpleQueueType = 0
	Transient SimpleQueueType = 1
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error{
	jsonBytes,err :=json.Marshal(val)
	if err!=nil{
		return fmt.Errorf("error %v", err)
	}
	ch.PublishWithContext(context.Background(),exchange,key,false,false,amqp.Publishing{ContentType: "application/json", Body: jsonBytes})
	fmt.Println("")
	return nil
}
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
    handler func(T),
) error{
	pubsub.DeclareAndBind(conn,exchange,queueName,key,queueType)
	// channel:= make(chan amqp.Delivery)
	// channel.Consume()
	channel,err:=conn.Channel()
	if err!=nil{
		return err
	}
	chDeli,err:=channel.Consume(queueName,"",false,false,false,false,nil)
	if err!=nil{
		return err
	}
	// go UnMarshal[T](chDeli)
	for val := range chDeli{
		var result []T
		json.Unmarshal(val.Body,result)
	}
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