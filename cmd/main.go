package main

import (
	"log"
	"rabbit-mq/util"
)

func main() {
	TestQueue()
}

func TestQueue() {
	inputData := util.SyncRequest{
		Topic: "123",
		Body:  "456",
	}

	var centerQueueConfig = util.QueueConfig{
		Host:     "localhost",
		Port:     "5672",
		Username: "guest",
		Password: "guest",
		PortUI:   "",
	}

	util.CreateExchange(centerQueueConfig, "test-exchange", "topic")
	centerConn := util.GetRabbitmqConn(centerQueueConfig)
	centerChannel := util.GetRabbitmqChannel(centerConn)

	err := util.CreateQueue("test-queue", true, centerChannel)
	if err != nil {
		log.Fatal(err)
	}
	util.BindQueue(centerQueueConfig, "test-exchange", "test", "test-queue")
	util.Produce(inputData, "test-exchange", "test", centerChannel)
	println("Push done")

	// block here until program is interrupted
	util.Consume("test-queue", centerChannel)
}
