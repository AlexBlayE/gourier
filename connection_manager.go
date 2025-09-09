package gourier

import (
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

type connectionManager struct {
	connections map[[4]byte]net.Conn
	mu          sync.Mutex

	deadlineTime uint

	maxBytes uint

	maxGoroutines uint
	grCount       uint
	grBlocker     chan<- struct{}

	radix *radixNode
}

func (cm *connectionManager) GetConn(ip [4]byte) net.Conn {
	cm.mu.Lock()
	conn, ok := cm.connections[ip]
	cm.mu.Unlock()
	if !ok {
		return nil
	}

	return conn
}

func (cm *connectionManager) SaveConn(ip [4]byte, conn net.Conn) {
	cm.mu.Lock()
	cm.connections[ip] = conn
	cm.mu.Unlock()
	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(cm.deadlineTime)))

	cm.grCount++ // TODO: fer que el count sigui amb sync
	// Ficar aquí el chan del blocker
	go func() {
		defer func() {
			cm.grCount--
			// cm.CloseConn(conn.LocalAddr().Network())
			// TODO: blocker chan lliberar
		}()

		for {
			buff := make([]byte, cm.maxBytes)
			lr := io.LimitReader(conn, int64(cm.maxBytes))

			_, err := lr.Read(buff)
			if err != nil {
				return
			}

			rn := cm.radix.FindPath(buff...)
			if rn == nil {
				return
			}

			hr := &handlerRunner{&Context{conn, cm.connections, buff, rn.depth}}
			hr.RunHandlers(rn.GetHandlers()...)
		}
	}()
}

func (cm *connectionManager) CloseConn(ip [4]byte) error {
	cm.mu.Lock()
	conn, ok := cm.connections[ip]
	if !ok {
		return errors.New("no existing connection")
	}
	delete(cm.connections, ip)
	cm.mu.Unlock()
	conn.Close()
	return nil
}
