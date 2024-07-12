package internal

import (
	"bytes"
	"encoding/json"
)

type messsageType uint16

// constant values for message types
const (
	Publish     messsageType = iota + 1 // normal publish message
	Subscribe                           // subscribe message
	Unsubscribe                         // unsubscribe message
	PingMessage                         // ping message
	PongMessage                         // pong message
	Imposter                            // unauthorized user message
)

// message is what we send between worker and clients.
type message struct {
	Type  messsageType `json:"type"`
	Topic string       `json:"topic,omitempty"`
	Data  []byte       `json:"data,omitempty"`
}

// NewMessage generates a new message type.
func newMessage(t messsageType, topic string, data []byte) message {
	return message{
		Type:  t,
		Topic: topic,
		Data:  data,
	}
}

// EncodeMessage will convert message to array of bytes.
func encodeMessage(m message) []byte {
	bytes, _ := json.Marshal(m)
	return append(bytes, '\n')
}

// DecodeMessage will convert array of bytes to Message.
func decodeMessages(buffer []byte) ([]message, error) {
	fields := bytes.Fields(buffer)

	messages := make([]message, 0, len(fields))

	for _, field := range fields {
		var message message
		if err := json.Unmarshal(field, &message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
