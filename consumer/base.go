package consumer

import (
	"github.com/streadway/amqp"
)

type Consume struct {
	Channel *amqp.Channel
}

// GetConn -
func GetConn(rabbitURL string) (Consume, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return Consume{}, err
	}

	ch, err := conn.Channel()
	return Consume{
		Channel: ch,
	}, err
}
