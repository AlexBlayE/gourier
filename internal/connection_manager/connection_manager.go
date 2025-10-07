package connectionmanager

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/AlexBlayE/gourier"
	"github.com/AlexBlayE/gourier/internal/context"
	handlerrunner "github.com/AlexBlayE/gourier/internal/handle_runner"
)

type ConnectionManager struct {
	deadlineTime uint

	maxBytes uint

	pathFinder gourier.PathFinder

	sPool *sync.Pool
}

func NewConnectionManager(
	deadlineTime uint,
	maxReadBytes uint,
	pathFinder gourier.PathFinder,
) *ConnectionManager {
	return &ConnectionManager{
		deadlineTime: deadlineTime,
		maxBytes:     maxReadBytes,
		pathFinder:   pathFinder,
		sPool: &sync.Pool{
			New: func() any {
				return &context.Context{}
			},
		},
	}
}

func (cm *ConnectionManager) ManageConn(conn net.Conn) {
	for {
		conn.SetDeadline(time.Now().Add(time.Second * time.Duration(cm.deadlineTime)))

		b, err := cm.ReadAll(conn)
		if err != nil {
			return
		}

		rn := cm.pathFinder.FindPath(b...)
		if rn == nil {
			return
		}

		ctx := cm.sPool.Get().(*context.Context)
		ctx.Conn = conn
		ctx.Payload = b
		ctx.Depth = rn.GetDepth()
		ctx.AbortFlag = false
		ctx.Store = make(map[string]any)

		hr := &handlerrunner.HandlerRunner{Ctx: ctx}

		err = hr.RunHandlers(rn.GetHandlers()...)
		cm.sPool.Put(ctx)
		if err != nil {
			return
		}
	}
}

func (cm *ConnectionManager) WriteAll(conn net.Conn, data []byte) error {
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

func (cm *ConnectionManager) ReadAll(conn net.Conn) ([]byte, error) {
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
