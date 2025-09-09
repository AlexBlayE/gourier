package gourier

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type logger struct {
	formatter fmt.Formatter
	writer    io.Writer
	mu        sync.Mutex
	level     string
}

func (l *logger) Format(f fmt.State, c rune) {
	timestamp := time.Now().Format("0000-00-00 00:00:00")
	switch c {
	case 'v':
		if f.Flag('+') {
			fmt.Fprintf(f, "[%s] Logger{Level:%s}", timestamp, l.level)
		} else {
			fmt.Fprintf(f, "[%s] %s", timestamp, l.level)
		}
	case 's':
		fmt.Fprintf(f, "%s", l.level)
	}
}

func (l *logger) Log(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.writer, "[%s] %s\n", time.Now().Format("15:04:05"), msg)
}
