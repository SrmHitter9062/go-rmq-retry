package main

import (
	"fmt"

	"github.com/SrmHitter9062/go-rmq-retry/consumer"
	"github.com/SrmHitter9062/go-rmq-retry/processor"
	"github.com/streadway/amqp"
)

func init() {

}

func handler(d amqp.Delivery) bool {
	if d.Body == nil {
		fmt.Println("Error ,no message received")
		return false
	} else {
		fmt.Println("Message received:", string(d.Body))
		return true
	}
}

func initConsumer() {
	conn, err := consumer.GetConn("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	// start consuming
	err = conn.AddListener("srm-exchange", "srm-queue", "srm-key", processor.Process)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	<-forever
}

func main() {
	initConsumer()
}
