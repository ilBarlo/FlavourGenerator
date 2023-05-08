package flavourgenerator

// sendMessage sends a message on the queue
func sendMessage(flavour Flavour, queueName string, url string) error {

	message, err := marshallJson(&flavour)
	if err != nil {
		return err
	}
	// Channel creation
	conn, ch, err := createChannel(url)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer ch.Close()

	// Code declaration
	err = declareQueue(ch, queueName)
	if err != nil {
		return err
	}

	// Publish Message
	err = publishMessage(ch, queueName, message)
	if err != nil {
		return err
	}

	return nil
}
