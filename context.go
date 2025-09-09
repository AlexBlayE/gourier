package gourier

import (
	"net"
)

type HandleFunc func(ctx *Context)

type Context struct {
	conn  net.Conn
	conns map[[4]byte]net.Conn

	payload []byte
	depth   uint
}

// func (c *Context) Abort()

// func (c *Context) Send(payload []byte, headers ...byte)

func (c *Context) GetConn() net.Conn {
	return c.conn
}

func (c *Context) GetPayload() []byte {
	return c.payload[c.depth:]
}
