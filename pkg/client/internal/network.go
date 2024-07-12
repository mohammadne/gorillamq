package internal

import (
	"fmt"
	"io"
	"net"
)

// network handles the tcp requests.
type network struct {
	connection net.Conn
}

// send data over tcp.
func (n *network) send(data []byte) error {
	if _, err := n.connection.Write(data); err != nil {
		return fmt.Errorf("failed to send: %v\n", err)
	}

	return nil
}

// get data from tcp.
func (n *network) get() ([]byte, error) {
	var buffer = make([]byte, 2048)
	bytes, err := n.connection.Read(buffer)
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		LogError("network read error", err)
	}

	return buffer[:bytes], nil
}
