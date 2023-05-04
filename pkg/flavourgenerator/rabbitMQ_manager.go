package flavourgenerator

import (
	"github.com/streadway/amqp"
)

// createChannel creates a RabbitMQ Channel
func createChannel(url string) (*amqp.Connection, *amqp.Channel, error) {
	// Connection to the server RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}
	// Create Channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	return conn, ch, nil
}

// declareQueue declares a queue on a RabbitMQ Channel
func declareQueue(ch *amqp.Channel, queueName string) error {

	_, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}
	return nil
}

// publishMessage publish a message on the channel
func publishMessage(ch *amqp.Channel, queueName string, message []byte) error {

	err := ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
