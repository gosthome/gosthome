package cv

import (
	"bytes"
	"context"
	"log/slog"
)

type Bytes struct {
	Data []byte
}

func (b *Bytes) ValidateWithContext(ctx context.Context) error {
	return nil
}
func (b *Bytes) Equal(other *Bytes) bool {
	return bytes.Equal(b.Data, other.Data)
}
func (b *Bytes) MarshalText(text []byte) (err error) {
	panic("unimpl")
}
func (b *Bytes) UnmarshalText(text []byte) (err error) {
	slog.Info("(b *Bytes) UnmarshalText(", "text", text)
	b.Data = text
	return nil
}
