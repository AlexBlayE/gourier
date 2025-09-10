package gourier

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"
)

type connectionManager struct {
	deadlineTime uint

	maxBytes uint

	grBlocker chan struct{}

	radix *radixNode
}

func newConnectionManager(deadlineTime uint, maxReadBytes uint, maxGoroutines uint, radix *radixNode) *connectionManager {
	return &connectionManager{
		deadlineTime: deadlineTime,
		maxBytes:     maxReadBytes,
		grBlocker:    make(chan struct{}, maxGoroutines),
		radix:        radix,
	}
}

func (cm *connectionManager) ManageConn(conn net.Conn) {
	cm.grBlocker <- struct{}{}
	go func() {
		defer func() {
			conn.Close()
			<-cm.grBlocker
		}()

		for {
			buff := make([]byte, cm.maxBytes)

			conn.SetDeadline(time.Now().Add(time.Second * time.Duration(cm.deadlineTime)))

			err := cm.readAll(conn, buff)
			if err != nil {
				return
			}

			rn := cm.radix.FindPath(buff...)
			if rn == nil {
				return
			}

			hr := &handlerRunner{&Context{conn, buff, rn.depth, false}}

			err = hr.RunHandlers(rn.GetHandlers()...)
			if err != nil {
				return
			}
		}
	}()
}

func (cm *connectionManager) writeAll(conn net.Conn, data []byte) error {
	total := 0

	for total < len(data) {
		n, err := conn.Write(data[total:])
		if err != nil {
			return err
		}
		total += n
	}

	return nil
}

func (cm *connectionManager) readAll(conn net.Conn, dst []byte) error {
	encodedLength := [4]byte{}

	_, err := io.ReadFull(conn, encodedLength[:])
	if err != nil {
		return err
	}

	payloadLength := binary.BigEndian.Uint32(encodedLength[:])
	if payloadLength > uint32(cm.maxBytes) {
		return errors.New("payload too large")
	}

	_, err = io.ReadFull(conn, dst[:payloadLength])
	if err != nil {
		return err
	}

	return nil
}
