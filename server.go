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
	radixRouter *radixNode
}

func New() *server {
	radix := &radixNode{make(map[byte]*radixNode), nil, nil, 0}

	return &server{
		connManager: newConnectionManager(60, 1024, 50, radix),
		logger:      &logger{nil, os.Stdout, sync.Mutex{}, "Info"},
		radixRouter: radix,
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

		go s.connManager.ManageConn(conn)
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

		go s.connManager.ManageConn(conn)
	}
}

func (s *server) Handler(header byte, handleFunc ...HandleFunc) {
	s.radixRouter.children[header] = &radixNode{nil, handleFunc, nil, 1}
}

func (s *server) Group(header byte) *routerGroup {
	newChildNode := &radixNode{make(map[byte]*radixNode), nil, nil, 1}
	s.radixRouter.children[header] = newChildNode
	return &routerGroup{newChildNode}
}

func (s *server) SetOptions() error {
	// TODO: opcions del connManager com MaxGoroutines, MaxBytes etc
	return nil
}

func (s *server) HandleOpenConn(conn net.Conn) {
	s.connManager.ManageConn(conn)
}

func (s *server) Send(ip string) {

}

func (s *server) SendTLS(ip string, payload []byte, config *tls.Config) error {
	conn, err := tls.Dial("tcp", ip, config)
	if err != nil {
		return err
	}

	// conn.Write()// TODO: fer que s'envii tot el payload mirant n(crear una funció que ho fagi automaticament)
	// TODO: Com que tcp es orientet al fluxe tinc que fer un mecanisme perque send y write sempre ho llegeixin tot y despres ya continui la execució

	s.HandleOpenConn(conn)
	return nil
}
