package main

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.Acktype {
	return func(ps routing.PlayingState) pubsub.Acktype{
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerMove(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.ArmyMove) pubsub.Acktype {
	var moveOutcome gamelogic.MoveOutcome
	return func(mv gamelogic.ArmyMove ) pubsub.Acktype{
		defer fmt.Print("> ")
		moveOutcome=gs.HandleMove(mv)
		switch moveOutcome{
			case gamelogic.MoveOutComeSafe:
				return pubsub.Ack
			case gamelogic.MoveOutcomeMakeWar:
				if err:= pubsub.PublishJSON(publishCh,routing.ExchangePerilTopic,routing.WarRecognitionsPrefix+"."+gs.GetPlayerSnap().Username,gamelogic.RecognitionOfWar{
					Attacker: mv.Player,
					Defender: gs.GetPlayerSnap(),
				});err!=nil{
					fmt.Println("failed to publish war message:", err)
					return pubsub.NackRequeue
				}
				fmt.Print("Sending Acknowledgement: Processed successfully.")
				return pubsub.Ack
			case gamelogic.MoveOutcomeSamePlayer:
				return pubsub.NackDiscar
			default:
				return pubsub.NackDiscar
		}
	}
}

func handlerWar(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(rw gamelogic.RecognitionOfWar) pubsub.Acktype{
		defer fmt.Print("> ")
		moveOutcome, winner, loser := gs.HandleWar(rw)
		switch moveOutcome {
		case gamelogic.WarOutcomeNotInvolved:
			fmt.Print("Sending Nack and requeue: Not processed successfully, WarOutcomeNotInvolved (retry).")
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			fmt.Print("Sending Nack and discard: Not processed successfully, and should be discarded (to a dead-letter queue if configured or just deleted entirely).")
			return pubsub.NackDiscar
		case gamelogic.WarOutcomeOpponentWon:
			msg := fmt.Sprintf("%s won a war against %s", winner, loser)
			if err := publishGameLog(publishCh, rw.Attacker.Username, msg); err != nil {
				fmt.Println("failed to publish game log:", err)
				return pubsub.NackRequeue
			}
			fmt.Print("Sending Acknowledgement: Processed successfully.")
			return pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			msg := fmt.Sprintf("%s won a war against %s", winner, loser)
			if err := publishGameLog(publishCh, rw.Attacker.Username, msg); err != nil {
				fmt.Println("failed to publish game log:", err)
				return pubsub.NackRequeue
			}
			fmt.Print("Sending Acknowledgement: Processed successfully.")
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			msg := fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			if err := publishGameLog(publishCh, rw.Attacker.Username, msg); err != nil {
				fmt.Println("failed to publish game log:", err)
				return pubsub.NackRequeue
			}
			fmt.Print("Sending Acknowledgement: Processed successfully.")
			return pubsub.Ack
		default:
			fmt.Print("Sending Nack and discard: Not processed successfully, and should be discarded (to a dead-letter queue if configured or just deleted entirely).")
			return pubsub.NackDiscar
		}

	}
}

func publishGameLog(ch *amqp.Channel, username string, msg string) error{
	gl:= routing.GameLog{
		CurrentTime: time.Now(),
		Message: msg,
		Username: username,
	}
	return pubsub.PublishGob(ch,string(routing.ExchangePerilTopic),routing.GameLogSlug+"."+username, gl)
}