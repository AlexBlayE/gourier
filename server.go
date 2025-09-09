package gourier

import (
	"crypto/tls"
	"net"
	"os"
	"sync"
)

type server struct {
	connManager *connectionManager
	logger      *logger
	routes      *radixNode
}

func New() *server {
	return &server{
		connManager: &connectionManager{},
		logger:      &logger{nil, os.Stdout, sync.Mutex{}, "Info"},
		routes:      &radixNode{make(map[byte]*radixNode), nil, 0},
	}
}

func (s *server) Run(port string) error {
	l, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *server) RunTLS(port string, tlsConfig *tls.Config) error {
	l, err := tls.Listen("tcp", port, tlsConfig)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *server) Handler(header byte, handleFunc ...HandleFunc) {
	s.routes.children[header] = &radixNode{nil, handleFunc, 1}
}

func (s *server) Group(header byte) *routerGroup {
	newChildNode := &radixNode{make(map[byte]*radixNode), nil, 1}
	s.routes.children[header] = newChildNode
	return &routerGroup{newChildNode}
}

func (s *server) SetOptions() error {
	// TODO: permetre opcions de gestio del tamañ de la petició
	// TODO: fer que les conexions en goroutines estinguin dins un pool de gooroutines per controlar les conexións maximes simultanies
	return nil
}

func (s *server) handleConnection(conn net.Conn) {
	// s.connManager.SaveConn(nil, conn) // TODO:
}
