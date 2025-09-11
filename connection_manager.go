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
			conn.SetDeadline(time.Now().Add(time.Second * time.Duration(cm.deadlineTime)))

			b, err := cm.readAll(conn)
			if err != nil {
				return
			}

			rn := cm.radix.FindPath(b...)
			if rn == nil {
				return
			}

			hr := &handlerRunner{&Context{conn, b, rn.depth, false, make(map[string]any)}}

			err = hr.RunHandlers(rn.GetHandlers()...)
			if err != nil {
				return
			}
		}
	}()
}

func (cm *connectionManager) writeAll(conn net.Conn, data []byte) error {
	size := len(data)
	lengthHeader := []byte{}

	lengthHeader = binary.BigEndian.AppendUint32(lengthHeader, uint32(size))

	totalHeader := 0
	for totalHeader < 4 {
		n, err := conn.Write(lengthHeader[totalHeader:])
		if err != nil {
			return err
		}
		totalHeader += n
	}

	total := 0
	for total < size {
		n, err := conn.Write(data[total:])
		if err != nil {
			return err
		}
		total += n
	}

	return nil
}

func (cm *connectionManager) readAll(conn net.Conn) ([]byte, error) {
	encodedLength := [4]byte{}

	_, err := io.ReadFull(conn, encodedLength[:])
	if err != nil {
		return nil, err
	}

	payloadLength := binary.BigEndian.Uint32(encodedLength[:])
	if payloadLength > uint32(cm.maxBytes) {
		return nil, errors.New("payload too large")
	}

	var b []byte = make([]byte, payloadLength)
	_, err = io.ReadFull(conn, b[:payloadLength])
	if err != nil {
		return nil, err
	}

	return b, nil
}
