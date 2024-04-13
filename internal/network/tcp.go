package network

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
)

// func ListenInsecureTCP(address string, handler func(Network)) error {
// 	listener, err := net.Listen("tcp", address)
// 	if err != nil {
// 		return fmt.Errorf("failed to start the server: %v", err)
// 	}

// 	for {
// 		if conn, er := listener.Accept(); er == nil {
// 			go handler(&tcp{connection: conn})
// 		}
// 	}
// }

// Listen will listens and serve
func ListenTCP(cfg *Config, address string, handler func(Network)) error {
	// Load server certificate and private key
	cert, err := tls.X509KeyPair([]byte(cfg.TLS.Certificate), []byte(cfg.TLS.PrivateKey))
	if err != nil {
		return fmt.Errorf("Error loading certificate: %v", err)
	}

	// Create TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Create TLS listener
	listener, err := tls.Listen("tcp", address, tlsConfig)
	if err != nil {
		return fmt.Errorf("Error creating listener: %v", err)
	}
	defer listener.Close()

	for {
		if conn, err := listener.Accept(); err == nil {
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
