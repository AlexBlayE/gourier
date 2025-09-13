package gourier

import (
	"crypto/tls"
	"net"
)

type Server struct {
	connManager  *connectionManager // TODO: que el connection manager sigui una interface que tingui la funció ManageConn
	radixRouter  *radixNode         // TODO: que sigui una interface que tingui la funció FindPath (nose si pot ser perque utilitzo el s.radixRouter.depth directament)
	*routerGroup                    // TODO: cambiar per interfaz
}

func New() *Server {
	radix := &radixNode{make(map[byte]*radixNode), nil, nil, 0}

	return &Server{
		connManager: newConnectionManager(60, 1024, 50, radix),
		radixRouter: radix,
		routerGroup: &routerGroup{radix},
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

func (s *Server) SetOptions() error {
	// TODO: opcio ns del connManager com MaxGoroutines, MaxBytes etc
	return nil
}

func (s *Server) Send(ip string, payload []byte) error {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		return err
	}

	err = s.connManager.WriteAll(conn, payload)
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

	err = s.connManager.WriteAll(conn, payload)
	if err != nil {
		return err
	}

	s.connManager.ManageConn(conn)

	return nil
}
