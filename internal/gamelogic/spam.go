package gamelogic

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)
func (gs *GameState) CommandSpam(words []string,ch *amqp.Channel) error {
	
	if len(words) < 2 {
		return errors.New("usage: spam <n>")
	}
	num, err := strconv.Atoi(words[1])
	if err != nil {
		return fmt.Errorf("error: %s is not a valid number", words[1])
	}
	for i :=0;i<num;i++{
		gl:= routing.GameLog{
			CurrentTime: time.Now(),
			Message: GetMaliciousLog(),
			Username: gs.Player.Username,
		}
		err=pubsub.PublishGob(ch,routing.ExchangePerilTopic,routing.GameLogSlug+"."+gs.Player.Username,gl)
		if err!=nil{
			fmt.Printf("%s", err)
		}
	}
	return nil
}
