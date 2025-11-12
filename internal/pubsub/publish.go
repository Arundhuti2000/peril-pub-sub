package pubsub

import (
	"context"
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