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
	queueName  string
}

func NewRabbitMQ(connection *amqp.Connection, queueName string) (Notifier, error) {
	return rabbitMQ{
		connection: connection,
		queueName:  queueName,
	}, nil
}

//TODO: add error handling/logging!
func (r rabbitMQ) Notify(domain, request string) {
	channel, _ := r.connection.Channel()
	defer channel.Close()

	q, _ := channel.QueueDeclare(
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

	_ = channel.Publish(
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
