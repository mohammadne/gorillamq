package manager

import "github.com/mohammadne/gorillamq/internal/core"

type SubscribeEvent struct {
	topic    core.Topic
	clientID ClientID
	deliver  chan<- DeliverEvent
}

type UnsubscribeEvent struct {
	topic    core.Topic
	clientID ClientID
}

type PublishEvent struct {
	topic   core.Topic
	message *core.Message
}

type DeliverEvent struct {
	message *core.Message
}

// TerminateEvent sends an terminate event to the broker to terminate this client
// ie it means to remove this client from all the subscribers
type TerminateEvent struct {
	clientID ClientID
}
