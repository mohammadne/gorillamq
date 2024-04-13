package network

import "errors"

var (
	ErrorConnectionClosed = errors.New("")
)

type Network interface {
	// Send sends data over defined protocol
	Send(data []byte) error

	// Recieve recieves data from defined protocol
	Recieve() ([]byte, error)
}
