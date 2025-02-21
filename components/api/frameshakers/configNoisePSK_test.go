package frameshakers_test

import (
	"testing"

	"github.com/gosthome/gosthome/components/api/frameshakers"
	"github.com/matryer/is"
)

func TestNoisePSKRound(t *testing.T) {
	for range 100 {
		is := is.New(t)
		g, err := frameshakers.GenerateEncryptionKey()
		is.NoErr(err)
		b, err := g.MarshalText()
		is.NoErr(err)
		parsed := &frameshakers.ConfigNoisePSK{}
		err = parsed.UnmarshalText(b)
		is.NoErr(err)
		is.Equal(g, parsed)

	}
}
