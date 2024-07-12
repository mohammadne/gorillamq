package internal

import (
	"fmt"
	"net"
	"sync"
)

// we need safety timeout, to prevent send more than
// one message in http request.
const (
	safetyTimeout = 1
)

// MessageHandler is a handler for messages that come from subscribing.
type MessageHandler func([]byte)

// client is our user application handler.
type client struct {
	locker *sync.RWMutex

	// map of topics
	topics map[string]MessageHandler

	// communication channel allows a client to make
	// a connection channel between read data and subscribers
	communicateChannel chan message

	// terminate channel is used to close a subscribe channel
	terminateChannel chan int

	// network handles the client socket data transfers
	network network
}

// NewClient creates a new client handler.
func NewClient(conn net.Conn, auth string) (*client, error) {
	c := &client{
		locker:             &sync.RWMutex{},
		topics:             make(map[string]MessageHandler),
		communicateChannel: make(chan message),
		terminateChannel:   make(chan int),
		network: network{
			connection: conn,
		},
	}

	// send the ping message
	if err := c.ping([]byte(auth)); err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	// starting data reader
	go c.readDataFromServer()

	// start listening on channels
	go c.listen()

	return c, nil
}

// readDataFromServer gets all data from server.
func (c *client) readDataFromServer() {
	for {
		// read data from network
		packet, er := c.network.get()
		if er != nil {
			LogError("failed read data", er)
			break
		}

		messages, err := decodeMessages(packet)
		if err != nil {
			LogInfo(string(packet))
			LogError("failed in message parse", err)
		}

		for _, message := range messages {
			// decode message
			if message.Type == Publish {
				c.communicateChannel <- message
			}
		}
	}

	// close
	c.terminateChannel <- 1
}

// listen method watches channels for input data.
func (c *client) listen() {
	for {
		select {
		case data := <-c.communicateChannel:
			c.handle(data)
		case <-c.terminateChannel:
			c.close()
		}
	}
}

// handle will execute the topic handler method.
func (c *client) handle(m message) {
	if handler, ok := c.topics[m.Topic]; ok {
		handler(m.Data)
	}
}

// close will terminate everything.
func (c *client) close() {
	_ = c.network.connection.Close()
}

// send a ping message to gorillamq server.
func (c *client) ping(data []byte) error {
	// sending ping data as a message
	if err := c.network.send(encodeMessage(newMessage(PingMessage, "", data))); err != nil {
		return fmt.Errorf("failed to ping server: %w", err)
	}

	// read data from network
	packet, er := c.network.get()
	if er != nil {
		return fmt.Errorf("server failed to pong: %w", er)
	}

	// check for message
	messages, err := decodeMessages(packet)
	if err != nil {
		return fmt.Errorf("decode message failed")
	}

	for _, message := range messages {
		switch message.Type {
		case PongMessage:
			return nil
		case Imposter:
			return fmt.Errorf("unauthorized user")
		default:
			return fmt.Errorf("connection failed")
		}
	}

	return nil
}

// Publish will send a message to broker server.
func (c *client) Publish(topic string, data []byte) error {
	err := c.network.send(encodeMessage(newMessage(Publish, topic, data)))
	if err != nil {
		return err
	}
	return nil
}

// Subscribe subscribes over broker.
func (c *client) Subscribe(topic string, handler MessageHandler) {
	c.locker.Lock()
	c.topics[topic] = handler // set a handler for given topic
	c.locker.Unlock()

	// send an http request to broker server
	err := c.network.send(encodeMessage(newMessage(Subscribe, topic, nil)))
	if err != nil {
		LogError("failed send message to broker server", err)
	}
}

// Unsubscribe removes client from subscribing over a topic.
func (c *client) Unsubscribe(topic string) {
	_ = c.network.send(encodeMessage(newMessage(Unsubscribe, topic, nil)))

	delete(c.topics, topic) // remove topic and its handler
}
