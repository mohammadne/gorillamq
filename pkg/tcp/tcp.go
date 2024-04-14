package tcp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
)

type TCP interface {
	ListenInsecureTCP(handler HandleTCP) error
	ListenSecureTCP(handler HandleTCP) error
}

type HandleTCP func(net.Conn)

func NewTCP(cfg *Config) TCP {
	return &tcp{config: cfg}
}

type tcp struct {
	config *Config
}

func (tcp *tcp) ListenInsecureTCP(handler HandleTCP) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", tcp.config.InsecurePort))
	if err != nil {
		return fmt.Errorf("failed to start the server: %v", err)
	}

	for {
		if connection, er := listener.Accept(); er == nil {
			go handler(connection)
		}
	}
}

// Listen will listens and serve
func (tcp *tcp) ListenSecureTCP(handler HandleTCP) error {
	// Load server certificate and private key
	cert, err := tls.X509KeyPair([]byte(tcp.config.TLS.Certificate), []byte(tcp.config.TLS.PrivateKey))
	if err != nil {
		return fmt.Errorf("Error loading certificate: %v", err)
	}

	// Create TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Create TLS listener
	listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", tcp.config.SecurePort), tlsConfig)
	if err != nil {
		return fmt.Errorf("Error creating listener: %v", err)
	}
	defer listener.Close()

	for {
		if connection, err := listener.Accept(); err == nil {
			go handler(connection)
		}
	}
}

func Send(connection net.Conn, data []byte) error {
	if _, err := connection.Write(data); err != nil {
		return fmt.Errorf("failed to send: %v\n", err)
	}

	return nil
}

var (
	ErrorConnectionClosed = errors.New("ErrorConnectionClosed")
)

func Recieve(connection net.Conn) ([]byte, error) {
	var buffer = make([]byte, 2048)
	bytes, err := connection.Read(buffer)
	if err != nil {
		if err == io.EOF {
			return nil, ErrorConnectionClosed
		}
		return nil, err
	}
	return buffer[:bytes], nil
}
