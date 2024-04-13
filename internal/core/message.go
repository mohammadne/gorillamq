package core

import (
	"bytes"
	"encoding/json"
)

type MessageType uint16

// constant values for message types incoming from the clients
const (
	Publish     MessageType = iota + 1 // normal publish message
	Subscribe                          // subscribe message
	Unsubscribe                        // unsubscribe message
	PingMessage                        // ping message
	PongMessage                        // pong message
	Imposter                           // unauthorized user message
)

// Message is what we send between worker and clients.
type Message struct {
	Type  MessageType `json:"type"`
	Topic Topic       `json:"topic,omitempty"`
	Data  []byte      `json:"data,omitempty"`
}

// EncodeMessage will convert message to array of bytes.
func (m *Message) Encode() []byte {
	bytes, _ := json.Marshal(m)
	return append(bytes, '\n')
}

// DecodeMessages will convert array of bytes to Message.
func DecodeMessages(packet []byte) ([]Message, error) {
	fields := bytes.Fields(packet)
	messages := make([]Message, 0, len(fields))
	for _, field := range fields {
		var message Message
		if err := json.Unmarshal(field, &message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
