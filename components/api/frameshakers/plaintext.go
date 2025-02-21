package frameshakers

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

// readVarUint is similar to binary.Uvarint() but reads one byte at a time.
func readVarUint(r io.Reader) (uint64, error) {
	var buf [1]byte
	var x uint64
	var s uint
	for i := 0; ; i++ {
		if _, err := r.Read(buf[:]); err != nil {
			return 0, err
		}
		b := buf[0]
		if b < 0x80 {
			if i >= binary.MaxVarintLen64 || i == binary.MaxVarintLen64-1 && b > 1 {
				return 0, errors.New("overflow")
			}
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
}

func PlaintextServer(
	ctx context.Context,
	r io.Reader,
	w io.Writer,
	framer ServerFramer,
) (
	err error,
) {
	msgData := make([]byte, 4096)
	var handler FramesHandler
	defer func() {
		if handler != nil {
			handler.Close()
		}
	}()

	readPacket := func() (int, []byte, error) {
		// read hello header
		h, rerr := readVarUint(r)
		if rerr != nil {
			return 0, nil, rerr
		}
		if h != 0x0 {
			return 0, nil, fmt.Errorf("header marker byte for plaintext invalid: %x", h)
		}
		msgLen, rerr := readVarUint(r)
		if rerr != nil {
			return 0, nil, rerr
		}
		type_, rerr := readVarUint(r)
		if rerr != nil {
			return 0, nil, rerr
		}
		if msgLen != 0 {
			msgData = reserveBuf(msgData, int(msgLen))
			msgData = msgData[:msgLen]
			_, rerr = r.Read(msgData)
			if rerr != nil {
				return 0, nil, rerr
			}
		} else {
			msgData = msgData[:0]
		}
		return int(type_), msgData, nil
	}
	writeBufMux := sync.Mutex{}
	writeBuf := make([]byte, 4096)
	writePacket := func(type_ int, packet []byte) error {
		writeBufMux.Lock()
		defer writeBufMux.Unlock()
		maxFrameLen := len(packet) + binary.MaxVarintLen64*2 + 1
		writeBuf = reserveBuf(writeBuf, maxFrameLen)
		writeBuf = writeBuf[:maxFrameLen]
		writeBuf[0] = 0
		p := 1
		p += binary.PutUvarint(writeBuf[p:], uint64(len(packet)))
		p += binary.PutUvarint(writeBuf[p:], uint64(type_))
		p += copy(writeBuf[p:], packet)
		writeBuf = writeBuf[:p]
		slog.Debug("Sending frame", "maxFrameLen", maxFrameLen, "written", p)
		if _, werr := w.Write(writeBuf); werr != nil {
			return werr
		}
		return nil
	}
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		type_, data, err := readPacket()
		if err != nil {
			return err
		}
		if handler == nil {
			handler, err = framer(func(frames []Frame) error {
				for _, frame := range frames {
					werr := writePacket(frame.Type, frame.Data)
					if werr != nil {
						return werr
					}
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to init framer %w", err)
			}
		}
		closing := false
		wp, err := handler.Handle(ctx, []Frame{{Type: type_, Data: data}})
		if err != nil {
			if !errors.Is(err, ErrCloseConnection) {
				return fmt.Errorf("handler errored: %w", err)
			}
			closing = true
		}
		for _, frame := range wp {
			err = writePacket(frame.Type, frame.Data)
			if err != nil {
				return err
			}
		}
		if closing {
			return ErrCloseConnection
		}
	}
}
