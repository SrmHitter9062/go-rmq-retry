package consumer

import (
	"fmt"
	"os"
)

type consumerInterface interface {
	connect(path string) error
}

// AddListener adds consumer
func (c *Consume) AddListener(exchangeName, queueName, routingKey string, handler func(data []byte) error) (err error) {
	fmt.Println("Declaring exchange")
	err = c.Channel.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	c.DeclareDLX()
	fmt.Println("Declaring queue")
	// define dlx for the main queue
	args := map[string]interface{}{
		"x-dead-letter-exchange": "srm-exchange-dlx",
		//"x-dead-letter-exchange":    "",
		//"x-dead-letter-routing-key": "srm-queue-dlx",
	}
	_, err = c.Channel.QueueDeclare(queueName, true, false, false, false, args)
	if err != nil {
		return err
	}
	fmt.Println("Binding queue")
	err = c.Channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		return err
	}
	fmt.Println("Starting Consuming")
	msgs, err := c.Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	// prefetch 4x as many messages as we can handle at once
	// prefetchCount := concurrency * 4
	// err = conn.Channel.Qos(prefetchCount, 0, false)
	// if err != nil {
	// 	return err
	// }

	fmt.Println("Running different thread for processing msg")
	// process in another thread
	go func() {
		for msg := range msgs {
			fmt.Println("message headers1", msg.Headers)
			fmt.Println("message headers", msg.Headers["x-death"].([]interface{}))
			// TODO: use count from msg.Headers["x-death"] array to decide that we have to discard or reject the message

			fmt.Println("Message received", string(msg.Body))
			err := handler(msg.Body)
			if err != nil {
				//msg.Ack(false)
				fmt.Println("Message rejected")
				msg.Nack(false, false)
			} else {
				fmt.Println("Message acked")
				msg.Ack(false)
			}
		}
		fmt.Println("Rabbit consumer closed - critical Error")
		os.Exit(1)

	}()

	return nil

}

// DeclareDLX creates the dlx exchange and queue
func (c *Consume) DeclareDLX() {
	fmt.Println("Declaring the dlx exchange")
	err := c.Channel.ExchangeDeclare(
		"srm-exchange-dlx", // name
		"direct",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Declaring the dlx queue")
	args := map[string]interface{}{
		"x-message-ttl":          1 * 30 * 1000,  // in ms
		"x-dead-letter-exchange": "srm-exchange", // define the primary exchange as dlx for this queue
	}
	_, err = c.Channel.QueueDeclare("srm-queue-dlx", true, false, false, false, args)
	if err != nil {
		fmt.Println("Error in declaring dlx queue", err)
	}
	fmt.Println("Binding dlx queue")
	err = c.Channel.QueueBind("srm-queue-dlx", "srm-key", "srm-exchange-dlx", false, nil)
	if err != nil {
		panic(err)
	}

}
