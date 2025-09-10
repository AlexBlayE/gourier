package gourier

import (
	"encoding/binary"
	"net"
)

type HandleFunc func(ctx *Context)

type Context struct {
	conn net.Conn

	payload []byte
	depth   uint

	abortFlag bool
}

func (c *Context) Abort(errorResponse []byte) {
	// TODO: fer que pugui enviar la resposta d'error
	c.abortFlag = true
}

func (c *Context) Send(payload []byte, headers ...byte) error {
	toSend := []byte{}

	fullLength := len(payload) + len(headers)

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
