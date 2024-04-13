package manager

import (
	"github.com/mohammadne/gorillamq/internal/core"
	"github.com/mohammadne/gorillamq/internal/network"
	"go.uber.org/zap"
)

type ClientID int

// Client represent each connection to the server
// which can subscribe on multiple topics and send
// messages over multiple topics
// Client handles a single client either for subscribing or publishing
type Client struct {
	logger  *zap.Logger
	network network.Network

	// the ID for each client routine
	ID ClientID

	subscribe   chan<- SubscribeEvent   // send subscribe event to the broker
	unsubscribe chan<- UnsubscribeEvent // send unsubscribe event to the broker
	publish     chan<- PublishEvent     // send publish event to the broker
	terminate   chan<- TerminateEvent   // send terminate event to the broker

	deliver chan DeliverEvent // recieve deliver event from the broker and pass it to the client
}

func NewClient(
	logger *zap.Logger,
	network network.Network,
	ID ClientID,
	subscribe chan<- SubscribeEvent,
	unsubscribe chan<- UnsubscribeEvent,
	publish chan<- PublishEvent,
	terminate chan<- TerminateEvent,
) *Client {
	return &Client{
		logger:      logger,
		network:     network,
		ID:          ID,
		subscribe:   subscribe,
		unsubscribe: unsubscribe,
		publish:     publish,
		terminate:   terminate,
		deliver:     make(chan DeliverEvent),
	}
}

func (c *Client) start() {
	c.logger.Info("client has been connected", zap.Int("ID", int(c.ID)))

	// handle handles incoming messages receiving from the network
	go func() {
		for {
			buffer, err := c.network.Recieve()
			if err != nil {
				if err == network.ErrorConnectionClosed {
					close(c.deliver)
					break
				}

				c.logger.Error("failed to read incoming network data", zap.Error(err))
				break
			}

			messages, err := core.DecodeMessages(buffer)
			if err != nil {
				c.logger.Error("failed to parse incoming network data", zap.ByteString("msg", buffer), zap.Error(err))
				break
			}

			for _, message := range messages {
				switch message.Type {
				case core.Publish:
					c.publish <- PublishEvent{
						topic:   message.Topic,
						message: &message,
					}
				case core.Subscribe:
					c.subscribe <- SubscribeEvent{
						topic:    message.Topic,
						clientID: c.ID,
						deliver:  c.deliver,
					}
				case core.Unsubscribe:
					c.unsubscribe <- UnsubscribeEvent{
						topic:    message.Topic,
						clientID: c.ID,
					}
				case core.PingMessage:
					msg := core.Message{Type: core.PongMessage}
					if err := c.network.Send(msg.Encode()); err != nil {
						c.logger.Error("Error sending pong message to the client", zap.Error(err))
					}
				default:
					c.logger.Error("unsupported or unimplemented message type has been given", zap.Any("message", message), zap.Int("ID", int(c.ID)))
				}
			}
		}

		// announcing that the worker is done
		c.terminate <- TerminateEvent{clientID: c.ID}
	}()

	// listen on deliver events and send the message to the client
loop:
	for {
		select {
		case event, more := <-c.deliver:
			if !more {
				break loop
			}
			if err := c.network.Send(event.message.Encode()); err != nil {
				c.logger.Error("Error sending message to the client", zap.Error(err))
			}
		}
	}
}
