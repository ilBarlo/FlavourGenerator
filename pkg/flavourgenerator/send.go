package flavourgenerator

import "log"

// sendMessage sends a message on the queue
func sendMessage(flavour Flavour, subject string, url string) error {

	message, err := marshallJson(&flavour)
	if err != nil {
		return err
	}

	// Connect to NATS server
	nc, err := connectNATS(url)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	// Publish a message to a subject
	if err := publishMsg(nc, subject, message); err != nil {
		log.Fatalf("Error publishing message: %v", err)
	}

	return nil
}
