package notifier

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

type Notifier interface {
	Notify(domain, request string)
}

type rabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
}

func NewRabbitMQ(connection *amqp.Connection, queueName string) (Notifier, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	return rabbitMQ{
		connection: connection,
		channel:    channel,
		queueName:  queueName,
	}, nil
}

func (r rabbitMQ) Notify(domain, request string) {
	q, _ := r.channel.QueueDeclare(
		r.queueName,
		false,
		false,
		false, false,
		nil,
	)

	message, _ := json.Marshal(struct {
		Domain  string `json:"domain"`
		Request string `json:"request"`
	}{
		Domain:  domain,
		Request: request,
	})

	_ = r.channel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}
