package flavourgenerator

// Function to send a message on the queue
func SendMessage(nodeInfo NodeInfo, queueName string, url string) error {

	message, err := marshallJson(&nodeInfo)
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
