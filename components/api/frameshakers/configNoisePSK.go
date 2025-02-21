package frameshakers

import (
	"bytes"
	"crypto/rand"
	"encoding"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"slices"
)

// tokenBytes generates a random byte slice of the specified length.
// If nbytes is 0, it returns a default number of bytes (e.g., 32).
func tokenBytes(nbytes int) ([]byte, error) {
	// Create a byte slice to hold the random bytes
	bytes := make([]byte, nbytes)

	// Read random bytes from the crypto/rand package
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func GenerateEncryptionKey() (n *ConfigNoisePSK, err error) {
	n = &ConfigNoisePSK{}
	// Use the tokenBytes function to generate random bytes
	n.data, err = tokenBytes(32)
	return
}

type ConfigNoisePSK struct {
	data []byte
}

func (n *ConfigNoisePSK) Equal(other *ConfigNoisePSK) bool {
	nv := n.Valid()
	ov := other.Valid()
	if nv && ov {
		return bytes.Equal(n.data, other.data)
	}
	return !nv && !ov
}

func (n *ConfigNoisePSK) Valid() bool {
	if n == nil {
		return false
	}
	if len(n.data) != 32 {
		slog.Error("noise psk key has wrong length", "actual", len(n.data))
		return false
	}
	return true
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *ConfigNoisePSK) UnmarshalText(text []byte) error {
	b := bytes.NewBuffer(make([]byte, 0, base64.StdEncoding.DecodedLen(len(text))))
	r := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(text))
	ld, err := io.Copy(b, r)
	if err != nil {
		return err
	}
	if ld != 32 {
		return fmt.Errorf("wrong decoded psk length (expected 32, got %d)", ld)
	}
	n.data = slices.Clone(b.Bytes())
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (n *ConfigNoisePSK) MarshalText() ([]byte, error) {
	b := bytes.NewBuffer(make([]byte, 0, base64.StdEncoding.EncodedLen(len(n.data))))
	w := base64.NewEncoder(base64.StdEncoding, b)
	_, err := w.Write(n.data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return slices.Clone(b.Bytes()), nil
}

func ParseNoisePSK(psk string) (*ConfigNoisePSK, error) {
	r := &ConfigNoisePSK{}
	err := r.UnmarshalText([]byte(psk))
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (n *ConfigNoisePSK) String() string {
	if n == nil {
		return ""
	}
	t, _ := n.MarshalText()
	return string(t)
}

func (n *ConfigNoisePSK) Data() []byte {
	if !n.Valid() {
		return nil
	}
	return n.data
}

var _ encoding.TextMarshaler = (*ConfigNoisePSK)(nil)
var _ encoding.TextUnmarshaler = (*ConfigNoisePSK)(nil)
