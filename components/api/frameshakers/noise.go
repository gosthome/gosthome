package frameshakers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"sync"

	"github.com/flynn/noise"
)

type handshakeError struct {
	msg []byte
	err error
}

func (he *handshakeError) Error() string {
	var add string
	if he.err != nil {
		add = he.err.Error()
	}
	return fmt.Sprintf("noise handshake: %s (%s)", string(he.msg), add)
}

func errBadIndicatorByte(err error) error {
	return &handshakeError{msg: []byte("Bad indicator byte"), err: err}
}
func errBadHandshakePacketLen(err error) error {
	return &handshakeError{msg: []byte("Bad handshake packet len"), err: err}
}
func errEmptyHandshakeMessage(err error) error {
	return &handshakeError{msg: []byte("Empty handshake message"), err: err}
}
func errBadHandshakeErrorByte(err error) error {
	return &handshakeError{msg: []byte("Bad handshake error byte"), err: err}
}
func errHandshakeMacFailure(err error) error {
	return &handshakeError{msg: []byte("Handshake MAC failure"), err: err}
}
func errHandshakeError(err error) error {
	return &handshakeError{msg: []byte("Handshake error"), err: err}
}

func getServerName(ctx context.Context) []byte {
	serverNameAny := ctx.Value(shakersKey{"serverName"})
	if serverNameAny == nil {
		return nil
	}
	serverName, ok := serverNameAny.(string)
	if !ok {
		return nil
	}
	return []byte(serverName)
}

func getNoisePSK(ctx context.Context) *ConfigNoisePSK {
	maybePSK := ctx.Value(shakersKey{"noisePSK"})
	if maybePSK == nil {
		return nil
	}
	psk, ok := maybePSK.(*ConfigNoisePSK)
	if !ok {
		return nil
	}
	return psk
}

type noiseStates int

const (
	noiseHello noiseStates = iota
	noiseHandshake
	noiseReady
)

func NoiseServer(
	ctx context.Context,
	r io.Reader,
	w io.Writer,
	framer ServerFramer,
) (
	err error,
) {
	serverName := getServerName(ctx)
	noisePSK := getNoisePSK(ctx)
	if !noisePSK.Valid() {
		return errors.New("invalid psk")
	}
	handshake, err := noise.NewHandshakeState(noise.Config{
		CipherSuite: noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashSHA256),

		Initiator: false,
		Prologue:  []byte("NoiseAPIInit\x00\x00"),
		Pattern:   noise.HandshakeNN,

		PresharedKey:          noisePSK.Data(),
		PresharedKeyPlacement: 0,
	})
	var enc, dec *noise.CipherState
	encMux := sync.Mutex{}
	encBuf := make([]byte, 4096)
	var handler FramesHandler
	defer func() {
		if handler != nil {
			herr := handler.Close()
			if herr != nil {
				slog.Error("error closing handler")
			}
		}
	}()

	reader := newNoisePacketReader(r)
	writePacket := newNoisePacketWriter(w).Write

	writeFrames := func(frames []Frame) error {
		for _, frame := range frames {
			ferr := locked(&encMux, func() error {
				data_len := len(frame.Data)
				encBuf = reserveBuf(encBuf, 4+data_len)
				encBuf[0] = byte((frame.Type >> 8) & 0xFF)
				encBuf[1] = byte(frame.Type & 0xFF)
				encBuf[2] = byte((data_len >> 8) & 0xFF)
				encBuf[3] = byte(data_len & 0xFF)
				encBuf = encBuf[:4+data_len]
				copy(encBuf[4:], frame.Data)
				slog.Debug("Encrypting message", "n", enc.Nonce(), "frameData", fmt.Sprintf("%x", frame.Data))
				sendBuf, eerr := enc.Encrypt(nil, nil, encBuf)
				if eerr != nil {
					return eerr
				}
				return writePacket(sendBuf)
			})
			if ferr != nil {
				return ferr
			}
		}
		return nil
	}
	maybeSendExplicitError := func(explErr error) error {
		if herr, ok := explErr.(*handshakeError); ok {
			err := writePacket([]byte{0x1}, herr.msg)
			if err != nil {
				return errors.Join(explErr, err)
			}
		}
		return explErr
	}
	state := noiseHello
	ctxDone := ctx.Done()
	var msgData []byte
	slog.Info("Entring noise encryption loop")
	for {
		select {
		case <-ctxDone:
			return ctx.Err()
		case fOrErr := <-reader.C:
			if fOrErr.err != nil {
				return maybeSendExplicitError(fOrErr.err)
			}
			msgData = fOrErr.data
		}

		if len(msgData) == 0 && state != noiseHello {

			return errors.New("received an empty message in non-hello state")
		}
		switch state {
		case noiseHello:
			err = writePacket([]byte{0x1}, serverName, []byte{0x0})
			if err != nil {
				return fmt.Errorf("failed to write hello %w", err)
			}
			state = noiseHandshake
		case noiseHandshake:
			if msgData[0] != 0x0 {
				return maybeSendExplicitError(errBadHandshakeErrorByte(errors.New("wrong marker delimiter for handshake")))
			}
			handshakeReply := make([]byte, 0)
			handshakeReply, _, _, err = handshake.ReadMessage(nil, msgData[1:])
			if err != nil {
				return maybeSendExplicitError(errHandshakeMacFailure(err))
			}
			handshakeReply, dec, enc, err = handshake.WriteMessage(handshakeReply[:0], nil)
			if err != nil {
				return maybeSendExplicitError(errHandshakeError(err))
			}
			if len(handshakeReply) == 0 {
				return errors.New("expected an established handshake, not nil")
			}
			err = writePacket([]byte{0x0}, handshakeReply)
			if err != nil {
				return fmt.Errorf("failed to write handshake: %w", err)
			}
			state = noiseReady
			handler, err = framer(writeFrames)
		case noiseReady:
			slog.Debug("Decrypting message", "n", dec.Nonce())
			msgData, err = dec.Decrypt(msgData[:0], nil, msgData)
			if err != nil {
				return err
			}
			msgType := (uint(msgData[0]) << 8) | uint(msgData[1])
			msgLen := (uint(msgData[2]) << 8) | uint(msgData[3])
			payload := msgData[4:]
			if len(payload) != int(msgLen) {
				return errors.New("message payload does not match sent length")
			}
			closing := false
			ret, err := handler.Handle(ctx, []Frame{{Type: int(msgType), Data: payload}})
			if err != nil {
				if !errors.Is(err, ErrCloseConnection) {
					return fmt.Errorf("handler errored: %w", err)
				}
				closing = true
			}
			err = writeFrames(ret)
			if err != nil {
				return fmt.Errorf("could not write frames %w", err)
			}
			if closing {
				return ErrCloseConnection
			}
		default:
			return errors.New("connection is in wrong state")
		}
	}
}

type noisePacketWriter struct {
	w     io.Writer
	bwMux sync.Mutex
	bw    *bytes.Buffer
}

func newNoisePacketWriter(w io.Writer) *noisePacketWriter {
	return &noisePacketWriter{
		w:     w,
		bwMux: sync.Mutex{},
		bw:    bytes.NewBuffer(make([]byte, 4096)),
	}
}

func (p *noisePacketWriter) Write(packetParts ...[]byte) error {
	p.bwMux.Lock()
	defer p.bwMux.Unlock()
	p.bw.Reset()
	var werr error
	if werr = p.bw.WriteByte(0x1); werr != nil {
		return werr
	}
	l := 0
	for _, part := range packetParts {
		l += len(part)
	}
	if werr = p.bw.WriteByte(byte((l >> 8) & 0xFF)); werr != nil {
		return werr
	}
	if werr = p.bw.WriteByte(byte(l & 0xFF)); werr != nil {
		return werr
	}
	for _, part := range packetParts {
		_, werr := p.bw.Write(part)
		if werr != nil {
			return werr
		}
	}
	slog.Debug("write frame", "data", fmt.Sprintf("%x", p.bw.Bytes()))
	if _, werr = p.w.Write(p.bw.Bytes()); werr != nil {
		return werr
	}
	return nil
}

type frameOrError struct {
	data []byte
	err  error
}

type noisePacketReader struct {
	sync.Mutex
	r io.Reader
	C <-chan frameOrError
	c chan<- frameOrError
}

func newNoisePacketReader(r io.Reader) *noisePacketReader {
	c := make(chan frameOrError, 1)
	ret := &noisePacketReader{
		Mutex: sync.Mutex{},
		r:     r,
		C:     c,
		c:     c,
	}
	go ret.run()
	return ret
}

func (pr *noisePacketReader) run() {
	// read hello header
	msgData := make([]byte, 4096)
	for {
		msgData = reserveBuf(msgData, 3)
		msgData = msgData[:3]
		n, rerr := pr.r.Read(msgData[:3])
		if rerr != nil {
			pr.c <- frameOrError{err: rerr}
			close(pr.c)
			return
		}
		if n != 3 {
			pr.c <- frameOrError{err: errors.New("short read during header")}
			close(pr.c)
			return
		}
		if msgData[0] != 0x1 {
			pr.c <- frameOrError{err: errBadIndicatorByte(fmt.Errorf("invalid: %x", msgData))}
			close(pr.c)
			return
		}
		msgLen := (uint(msgData[1]) << 8) | uint(msgData[2])
		slog.Debug("recieved header", "len", msgLen, "msg", fmt.Sprintf("%x", msgData[:3]))
		if msgLen != 0 {
			msgData = reserveBuf(msgData, int(msgLen))
			n, err := pr.r.Read(msgData[:msgLen])
			if err != nil {
				pr.c <- frameOrError{err: err}
				close(pr.c)
				return
			}
			msgData = msgData[:msgLen]
			if n != int(msgLen) {
				pr.c <- frameOrError{err: fmt.Errorf("short read while recieving message (expected %d, rscvd: %d)", n, msgLen)}
				close(pr.c)
				return
			}
			slog.Debug("recieved msg data", "msgLen", msgLen, "len(msgdata)", len(msgData))
		} else {
			msgData = msgData[:0]
		}
		pr.c <- frameOrError{data: slices.Clone(msgData)}
	}
}

func NoiseClient(
	ctx context.Context,
	r io.Reader,
	w io.Writer,
	framer ServerFramer,
) (
	err error,
) {

	// serverName := ""
	// noisePSK := getNoisePSK(ctx)
	// if !noisePSK.Valid() {
	// 	return errors.New("invalid psk")
	// }
	// handshake, err := noise.NewHandshakeState(noise.Config{
	// 	CipherSuite: noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashSHA256),

	// 	Initiator: true,
	// 	Prologue:  []byte("NoiseAPIInit\x00\x00"),
	// 	Pattern:   noise.HandshakeNN,

	// 	PresharedKey:          noisePSK.Data(),
	// 	PresharedKeyPlacement: 0,
	// })
	// var enc, dec *noise.CipherState

	return nil
}
