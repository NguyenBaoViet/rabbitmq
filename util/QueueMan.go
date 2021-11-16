package util

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type SyncRequest struct {
	Topic string
	Body  string
}
type QueueStats struct {
	Name         string
	MessageReady int64
	Consumers    int64
	IdleSince    int64
}

type QueueConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	PortUI   string
}

//exchange name
var JobExchange = "cr_ex_job"
var CrawlExchange = "cr_ex_crawl"
var LoadExchange = "cr_ex_load"
var TransExchange = "cr_ex_trans"

//prefix worker queues
var JobQueue = "cr_job"
var PrefixLoadQueue = "cr_load_"
var PrefixCrawlQueue = "cr_crawl_"
var TransQueue = "cr_trans"

//service global queue connection
var ServiceQueueConn *amqp.Connection
var ServiceQueueChannel *amqp.Channel

//create external queue queue channel instance
func GetRabbitmqConnChannel(config QueueConfig) (*amqp.Channel, *amqp.Connection) {
	conn, err := amqp.Dial("amqp://" + config.Username + ":" + config.Password + "@" + config.Host + ":" + config.Port + "/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ, Message: ", err)
	}
	//defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel, Message: ", err)
	}
	//defer ch.Close()

	//set prefetch
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatal("Failed to set QoS, Message: ", err)
	}

	return ch, conn
}

//all worker routine has same connection
func GetRabbitmqConn(config QueueConfig) *amqp.Connection {
	// rabbit mq config
	conn, err := amqp.Dial("amqp://" + config.Username + ":" + config.Password + "@" + config.Host + ":" + config.Port + "/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ, Message: ", err)
	}
	//defer conn.Close()

	return conn
}

//each worker routine has its own channel
func GetRabbitmqChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel, Message: ", err)
	}
	//defer ch.Close()

	//set prefetch
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatal("Failed to set QoS, Message: ", err)
	}
	return ch
}

//deployments queue architecture
func CreateExchange(config QueueConfig, exchange string, kind string) (err error) {
	ch, conn := GetRabbitmqConnChannel(config)
	defer ch.Close()
	defer conn.Close()

	//create exchange
	err = ch.ExchangeDeclare(
		exchange,
		kind,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare an exchange, Message: ", err)
	}
	return err
}

// create queue
func CreateQueue(queueName string, durabale bool, ch *amqp.Channel) (err error) {
	fmt.Print("yooo")

	_, err = ch.QueueDeclare(
		queueName,
		durabale,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// bind queue to exchange
func BindQueue(config QueueConfig, exchange string, routeKey string, queueName string) {
	ch, conn := GetRabbitmqConnChannel(config)
	defer ch.Close()
	defer conn.Close()
	err := ch.QueueBind(
		queueName,
		routeKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to bind a queue, message: ", err)
	}
}
