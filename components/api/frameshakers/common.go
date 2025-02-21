package frameshakers

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"runtime"
	"slices"
	"sync"
	"syscall"
)

type Frame struct {
	Type int
	Data []byte
}

type FramesHandler interface {
	Handle(ctx context.Context, input []Frame) ([]Frame, error)
	Close() error
}
type FrameSenderFunc = func([]Frame) error

type Shaker func(
	ctx context.Context,
	r io.Reader,
	w io.Writer,
	framer ServerFramer,
) (
	err error,
)

type ServerFramer func(sendFrames FrameSenderFunc) (handler FramesHandler, err error)

func SplitConnection(c net.Conn) (r io.Reader, w io.WriteCloser) {
	return bufio.NewReader(c), c
}

type shakersKey struct {
	Value string
}

func ContextWithValue(ctx context.Context, k string, v any) context.Context {
	return context.WithValue(ctx, shakersKey{k}, v)
}

func reserveBuf(b []byte, l int) []byte {
	lb := len(b)
	if lb < l {
		b = slices.Grow(b, l-lb)
	}
	return b
}

func locked(m *sync.Mutex, f func() error) error {
	m.Lock()
	defer m.Unlock()
	return f()
}

// isErrEOF returns true if the error is functionally equivalent to io.EOF.
//
// This is needed because the error is different on Windows.
func isErrEOF(err error) bool {
	if err == io.EOF {
		slog.Debug("isErrEOF(%T %#v) = true", "err", err)
		return true
	}
	if runtime.GOOS == "windows" {
		if oe, ok := err.(*net.OpError); ok && oe.Op == "read" {
			// Created by os.NewSyscallError()
			if se, ok := oe.Err.(*os.SyscallError); ok && se.Syscall == "wsarecv" {
				const WSAECONNABORTED = 10053
				const WSAECONNRESET = 10054
				switch n := se.Err.(type) {
				case syscall.Errno:
					v := n == WSAECONNRESET || n == WSAECONNABORTED
					slog.Debug("isErrEOF", "se.Err", se.Err, "v", v)
					return v
				default:
					slog.Debug("isErrEOF = false", "se.Err", se.Err)
					return false
				}
			}
		}
	}
	slog.Debug("isErrEOF(%T %#v) = false", "err", err)
	return false
}

var ErrCloseConnection = errors.New("the connection should be closed")
