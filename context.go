package gourier

import (
	"encoding/binary"
	"net"
)

type HandleFunc func(ctx *Context)

type Context struct { // TODO: fer que una interface abstraigui les seves funcions
	conn net.Conn

	payload []byte
	depth   uint

	abortFlag bool

	store map[string]any
}

func (c *Context) Abort(errorResponse []byte) {
	c.Send(errorResponse)
	c.abortFlag = true
}

func (c *Context) Send(payload []byte, headers ...byte) error {
	fullLength := len(payload) + len(headers)

	toSend := make([]byte, 4)

	binary.BigEndian.PutUint32(toSend, uint32(fullLength))

	for _, header := range headers {
		toSend = append(toSend, header)
	}

	toSend = append(toSend, payload...)

	total := 0
	for total < len(toSend) {
		n, err := c.conn.Write(toSend[total:])
		if err != nil {
			return err
		}
		total += n
	}

	return nil
}

func (c *Context) GetConn() net.Conn {
	return c.conn
}

func (c *Context) GetPayload() []byte {
	return c.payload[c.depth:]
}

func (c *Context) Get(key string) any {
	elem, ok := c.store[key]
	if !ok {
		return nil
	}

	return elem
}

func (c *Context) Set(key string, val any) {
	c.store[key] = val
}
