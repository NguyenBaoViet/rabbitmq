package util

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func Consume(inputQueue string, channel *amqp.Channel) {
	messages, err := channel.Consume(
		inputQueue, // queue
		"",         // consumer tag
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatal("Failed to register a consumer to " + inputQueue)
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			//Convert json data from queue
			var job interface{}
			jsonErr := json.Unmarshal(d.Body, &job)
			if jsonErr != nil {
				d.Ack(false)
				continue
			} else {
			}

			is_ok := true

			if is_ok == true {
				fmt.Println(job)
				d.Ack(false)
			} else {
				d.Nack(false, true)
			}
		} //get messages from queue loop
	}()

	fmt.Println()

	<-forever

}
