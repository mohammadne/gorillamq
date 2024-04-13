package manager

import (
	"context"
	"math/rand/v2"
	"sync"

	"github.com/mohammadne/gorillamq/internal/core"
	"github.com/mohammadne/gorillamq/internal/network"
	"go.uber.org/zap"
)

type broker struct {
	logger *zap.Logger
	locker *sync.RWMutex

	subscribe   chan SubscribeEvent   // recieve subscribe event from the client
	unsubscribe chan UnsubscribeEvent // recieve unsubscribe event from the client
	publish     chan PublishEvent     // recieve publish event from the client
	terminate   chan TerminateEvent   // recieve terminate event from the client

	listeners map[core.Topic]map[ClientID]chan<- DeliverEvent
}

const (
	subscribeChannelCapacity   = 5
	unsubscribeChannelCapacity = 3
	publishChannelCapacity     = 50
	terminateChannelCapacity   = 1
)

func NewBroker(logger *zap.Logger) *broker {
	return &broker{
		logger: logger,
		locker: &sync.RWMutex{},

		subscribe:   make(chan SubscribeEvent, subscribeChannelCapacity),
		unsubscribe: make(chan UnsubscribeEvent, unsubscribeChannelCapacity),
		publish:     make(chan PublishEvent, publishChannelCapacity),
		terminate:   make(chan TerminateEvent, terminateChannelCapacity),

		listeners: make(map[core.Topic]map[ClientID]chan<- DeliverEvent),
	}
}

func (b *broker) Start(ctx context.Context, cfg *network.Config) {
	err := network.ListenTCP(cfg, ":8080", func(network network.Network) {
		NewClient(
			b.logger,
			network,
			ClientID(rand.IntN(1000)), // TODO:change the logic
			b.subscribe,
			b.unsubscribe,
			b.publish,
			b.terminate,
		).start()
	})

	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case event := <-b.subscribe:
				if _, exists := b.listeners[event.topic]; !exists {
					b.listeners[event.topic] = map[ClientID]chan<- DeliverEvent{}
				}
				b.listeners[event.topic][event.clientID] = event.deliver
			case event := <-b.unsubscribe:
				delete(b.listeners[event.topic], event.clientID)
			case event := <-b.publish:
				topicListeners, exists := b.listeners[event.topic]
				if !exists {
					b.logger.Warn("the topic has no subscriber", zap.String("topic", string(event.topic)))
					continue
				}
				for _, listener := range topicListeners {
					listener <- DeliverEvent{message: event.message}
				}
			case event := <-b.terminate:
				for topic, topicListeners := range b.listeners {
					delete(topicListeners, event.clientID)
					if len(topicListeners) == 0 {
						delete(b.listeners, topic)
					}
				}
			}
		}
	}()

	// TODO
	// go func() {
	// 	for {
	// 		time.Sleep(time.Second)
	// 		b.logger.Info("broker information",
	// 			zap.Any("topics", b.Topics()),
	// 			zap.Any("len", len(b.listeners)),
	// 		)
	// 	}
	// }()
}

func (b *broker) Topics() []core.Topic {
	topics := make([]core.Topic, 0, len(b.listeners))
	for topic := range b.listeners {
		topics = append(topics, topic)
	}
	return topics
}
