package context

import (
	"encoding/binary"
	"net"
)

type Context struct {
	Conn net.Conn

	Payload []byte
	Depth   uint

	AbortFlag bool

	Store map[string]any
}

func (c *Context) SetAbortFlag(b bool) {
	c.AbortFlag = b
}

func (c *Context) GetAbortFlag() bool {
	return c.AbortFlag
}

func (c *Context) Abort(errorPayload []byte, headers ...byte) {
	c.Send(errorPayload, headers...)
	c.AbortFlag = true
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
		n, err := c.Conn.Write(toSend[total:])
		if err != nil {
			return err
		}
		total += n
	}

	return nil
}

func (c *Context) GetConn() net.Conn {
	return c.Conn
}

func (c *Context) GetPayload() []byte {
	return c.Payload[c.Depth:]
}

func (c *Context) Get(key string) any {
	elem, ok := c.Store[key]
	if !ok {
		return nil
	}

	return elem
}

func (c *Context) Set(key string, val any) {
	c.Store[key] = val
}
