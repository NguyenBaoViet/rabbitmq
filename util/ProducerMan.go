package util

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

func Produce(inputData interface{}, exchange string, routeKey string, ch *amqp.Channel) bool {
	isSuccess := false
	bodyJson, _ := json.Marshal(inputData)

	err := ch.Publish(
		exchange, // exchange
		routeKey, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bodyJson,
		})
	if err != nil {
		log.Fatal("Failed to publish a message to the queue, Message: ", err)
	}

	if err == nil {
		isSuccess = true
	}
	return isSuccess
}
