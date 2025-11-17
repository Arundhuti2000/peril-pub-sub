package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("error %v", err)
	}
	ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{ContentType: "application/json", Body: jsonBytes})
	fmt.Println("")
	return nil
}
func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(val)
	if err!=nil{
		return err
	}
	// jsonBytes, err := json.Marshal(val)
	if err!= nil {
		return fmt.Errorf("error %v", err)
	}
	ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{ContentType: "application/gob", Body: network.Bytes()})
	fmt.Println("")
	return nil
}