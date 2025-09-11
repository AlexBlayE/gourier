package gourier

import (
	"crypto/tls"
	"net"
	"os"
	"sync"
)

type Server struct {
	connManager *connectionManager
	logger      *logger
	radixRouter *radixNode
}

func New() *Server {
	radix := &radixNode{make(map[byte]*radixNode), nil, nil, 0}

	return &Server{
		connManager: newConnectionManager(60, 1024, 50, radix),
		logger:      &logger{nil, os.Stdout, sync.Mutex{}, "Info"},
		radixRouter: radix,
	}
}

func (s *Server) Run(port string) error {
	l, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go s.connManager.ManageConn(conn)
	}
}

func (s *Server) RunTLS(port string, tlsConfig *tls.Config) error {
	l, err := tls.Listen("tcp", port, tlsConfig)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go s.connManager.ManageConn(conn)
	}
}

func (s *Server) Handler(header byte, handleFunc ...HandleFunc) {
	s.radixRouter.children[header] = &radixNode{nil, handleFunc, nil, 1}
}

func (s *Server) Error(errorHandler HandleFunc) {
	s.radixRouter.errorHandler = errorHandler
}

func (s *Server) Group(header byte) *routerGroup {
	newChildNode := &radixNode{make(map[byte]*radixNode), nil, nil, 1}
	s.radixRouter.children[header] = newChildNode
	return &routerGroup{newChildNode}
}

func (s *Server) SetOptions() error {
	// TODO: opcions del connManager com MaxGoroutines, MaxBytes etc
	return nil
}

func (s *Server) Send(ip string, payload []byte) error {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		return err
	}

	err = s.connManager.writeAll(conn, payload)
	if err != nil {
		return err
	}

	s.connManager.ManageConn(conn)

	return nil
}

func (s *Server) SendTLS(ip string, payload []byte, config *tls.Config) error {
	conn, err := tls.Dial("tcp", ip, config)
	if err != nil {
		return err
	}

	err = s.connManager.writeAll(conn, payload)
	if err != nil {
		return err
	}

	s.connManager.ManageConn(conn)

	return nil
}
