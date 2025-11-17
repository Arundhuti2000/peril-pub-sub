package pubsub

import (
	"encoding/json"
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)
type SimpleQueueType int8
type Acktype int8


const(
	Durable SimpleQueueType = 0
	Transient SimpleQueueType = 1
)

const (
	Ack Acktype = 1
	NackRequeue Acktype = 2
	NackDiscar Acktype = 0
)

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) Acktype,
	unmarshaller func([]byte) (T, error),
) error {
	ch, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	chDeli, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("could not consume messages: %v", err)
	}

	go func() {
		defer ch.Close()
		for val := range chDeli {
			result, err := unmarshaller(val.Body)
			if err != nil {
				fmt.Printf("error unmarshalling: %v\n", err)
				val.Nack(false, true)
				continue
			}

			acktype := handler(result)

			switch acktype {
			case Ack:
				fmt.Print("Sending Acknowledgement: Processed successfully.")
				val.Ack(false)
			case NackRequeue:
				fmt.Print("Sending Nack and requeue: Not processed successfully, but should be requeued on the same queue to be processed again (retry).")
				val.Nack(false, true)
			case NackDiscar:
				fmt.Print("Sending Nack and discard: Not processed successfully, and should be discarded (to a dead-letter queue if configured or just deleted entirely).")
				val.Nack(false, false)
			default:
				fmt.Print("Error while decoding acknowledgment type")
			}
		}
	}()
	return nil
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
) error {
	unmarshaller := func(data []byte) (T, error) {
		var result T
		err := json.Unmarshal(data, &result)
		return result, err
	}
	return subscribe(conn, exchange, queueName, key, queueType, handler, unmarshaller)
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) Acktype,
	unmarshaller func([]byte) (T, error),
) error {
	return subscribe(conn, exchange, queueName, key, simpleQueueType, handler, unmarshaller)
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
		args:=amqp.Table{"x-dead-letter-exchange": routing.ExchangePerilDeadLetter}
		switch queueType {
			case 0:
				q,err= ch.QueueDeclare(queueName,true,false,false,false,args)
			case 1:
				q,err= ch.QueueDeclare(queueName,false, true,true, false, args)
			default:
				return ch, amqp.Queue{}, fmt.Errorf("could not find queue type: %v", err)
		}
		if err!=nil{
			return ch, amqp.Queue{}, fmt.Errorf("could not declare: %v", err)
		}
		err=ch.QueueBind(queueName,key,exchange,false,args)
		if err!=nil{
			return nil, amqp.Queue{},fmt.Errorf("could not bind queue: %v", err)
		}
		
		return ch, q, nil
	}