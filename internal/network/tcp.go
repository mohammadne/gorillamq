package network

import (
	"fmt"
	"io"
	"net"
)

func NewTCP() Network {
	return nil
}

// Listen will listens and serv
func ListenTCP(address string, handler func(Network)) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to start the server: %v", err)
	}

	for {
		if conn, er := listener.Accept(); er == nil {
			go handler(&tcp{connection: conn})
		}
	}
}

type tcp struct {
	connection net.Conn
}

func (tcp *tcp) Send(data []byte) error {
	if _, err := tcp.connection.Write(data); err != nil {
		return fmt.Errorf("failed to send: %v\n", err)
	}

	return nil
}

func (tcp *tcp) Recieve() ([]byte, error) {
	var buffer = make([]byte, 2048)
	bytes, err := tcp.connection.Read(buffer)
	if err != nil {
		if err == io.EOF {
			return nil, ErrorConnectionClosed
		}
		return nil, err
	}
	return buffer[:bytes], nil
}
