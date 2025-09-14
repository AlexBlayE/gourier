package gourier

import (
	"crypto/tls"
	"net"
)

type Context interface {
	GetConn() net.Conn
	GetPayload() []byte
	Abort(errorPayload []byte, headers ...byte)
	Send(payload []byte, headers ...byte) error
	Get(key string) any
	Set(key string, val any)
	SetAbortFlag(b bool)
	GetAbortFlag() bool
}

type HandleFunc func(ctx Context)

type ConnManager interface {
	ManageConn(conn net.Conn)
	ReadAll(conn net.Conn) ([]byte, error)
	WriteAll(conn net.Conn, data []byte) error
}

type PathFinder interface {
	FindPath(b ...byte) PathFinder
	GetDepth() uint
	GetHandlers() []HandleFunc
	SetErrorHandler(HandleFunc)
	SetChild(b byte, child PathFinder)
}

type RouteGroup interface {
	Error(errorHandler HandleFunc)
	Group(header byte) RouteGroup
	Handler(header byte, handleFunc ...HandleFunc)
}

type Server struct {
	connManager ConnManager

	RouteGroup

	gC chan struct{}
}

func New(connManager ConnManager, routeGroup RouteGroup, goroutinesChannel chan struct{}) *Server {
	return &Server{
		connManager: connManager,
		RouteGroup:  routeGroup,
		gC:          goroutinesChannel,
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

		s.gC <- struct{}{}
		go func() {
			defer func() {
				conn.Close()
				<-s.gC
			}()

			s.connManager.ManageConn(conn)
		}()
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

		s.gC <- struct{}{}
		go func() {
			defer func() {
				conn.Close()
				<-s.gC
			}()

			s.connManager.ManageConn(conn)
		}()
	}
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

// func Default() *Server {
// 	// radix := &pathfinder.RadixNode{make(map[byte]*pathfinder.RadixNode), nil, nil, 0}

// 	return &Server{
// 		// connManager: connectionmanager.NewConnectionManager(60, 1024, 50, radix),
// 		// pathFinder:  radix,
// 		// RouteGroup:  &group.RouterGroup{radix},
// 	}
// }
