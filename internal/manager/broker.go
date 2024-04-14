package manager

import (
	"context"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mohammadne/gorillamq/internal/core"
	"github.com/mohammadne/gorillamq/pkg/tcp"
	"go.uber.org/zap"
)

type broker struct {
	logger *zap.Logger
	tcp    tcp.TCP

	locker      *sync.RWMutex
	clientIdInc uint64

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

func NewBroker(logger *zap.Logger, tcp tcp.TCP) *broker {
	return &broker{
		logger: logger,
		tcp:    tcp,

		locker:      &sync.RWMutex{},
		clientIdInc: 1000000,

		subscribe:   make(chan SubscribeEvent, subscribeChannelCapacity),
		unsubscribe: make(chan UnsubscribeEvent, unsubscribeChannelCapacity),
		publish:     make(chan PublishEvent, publishChannelCapacity),
		terminate:   make(chan TerminateEvent, terminateChannelCapacity),

		listeners: make(map[core.Topic]map[ClientID]chan<- DeliverEvent),
	}
}

func (b *broker) Start(ctx context.Context) {
	connectionHandler := func(connection net.Conn) {
		NewClient(
			b.logger,
			connection,
			ClientID(b.clientIdInc),
			b.subscribe,
			b.unsubscribe,
			b.publish,
			b.terminate,
		).start()
		atomic.AddUint64(&b.clientIdInc, 1)
	}

	go func() {
		if err := b.tcp.ListenInsecureTCP(connectionHandler); err != nil {
			b.logger.Fatal("Error broker listen on insecure tcp connection", zap.Error(err))
		}
	}()

	go func() {
		if err := b.tcp.ListenSecureTCP(connectionHandler); err != nil {
			b.logger.Fatal("Error broker listen on secure tcp connection", zap.Error(err))
		}
	}()

	// send metrics about broker information
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case time := <-ticker.C:
				// TODO: enhance logic and metrics
				b.logger.Info("sending broker information",
					zap.Time("Time", time),
					// 			zap.Any("topics", b.Topics()),
					// 			zap.Any("len", len(b.listeners)),
				)
			case <-ctx.Done():
				return
			}
		}
	}()

	eventHandlers := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(eventHandlers)

	// swpan event handlers as standalone workers
	for index := 0; index < eventHandlers; index++ {
		go func(workerNumber int) {
			b.handleEvents(ctx, workerNumber)
			wg.Done()
		}(index + 1)
	}

	wg.Wait()
	b.logger.Info("gracefully existing the broker")
}

// topics returns a list of active topic which is listened to
func (b *broker) topics() []core.Topic {
	topics := make([]core.Topic, 0, len(b.listeners))
	for topic := range b.listeners {
		topics = append(topics, topic)
	}
	return topics
}

func (b *broker) handleEvents(ctx context.Context, workerIndex int) {
	b.logger.Info("worker start handling events", zap.Int("index", workerIndex))

loop:
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
		case <-ctx.Done():
			break loop
		}
	}
}
