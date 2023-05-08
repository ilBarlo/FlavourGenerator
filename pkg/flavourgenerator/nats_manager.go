package flavourgenerator

import (
	"github.com/nats-io/nats.go"
)

const natsURL = nats.DefaultURL

// connectNATS creates a connection to a NATS server
func connectNATS(url string) (*nats.Conn, error) {
	// Connect to NATS server
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return nc, nil
}

// publishMessage publishes a message to a subject
func publishMsg(nc *nats.Conn, subject string, message []byte) error {
	// Publish message
	err := nc.Publish(subject, message)
	if err != nil {
		return err
	}
	// Flush connection to ensure that message is sent
	err = nc.Flush()
	if err != nil {
		return err
	}
	return nil
}
