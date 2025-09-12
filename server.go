package gourier

import (
	"crypto/tls"
	"net"
)

// TODO: buscar utilitat e implementar interface que tingui Group, Error y Handle

type Server struct {
	connManager *connectionManager // TODO: que el connection manager sigui una interface que tingui la funció ManageConn
	radixRouter *radixNode         // TODO: que sigui una interface que tingui la funció FindPath
}

func New() *Server {
	radix := &radixNode{make(map[byte]*radixNode), nil, nil, 0}

	return &Server{
		connManager: newConnectionManager(60, 1024, 50, radix),
		radixRouter: radix,
	}
}

func (s *Server) Run(port string) error {
	l, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept() // TODO: ara mateix accepta totes les conexions aunque despues en ManageConn les bloqueji
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
	// TODO: opcio ns del connManager com MaxGoroutines, MaxBytes etc
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
