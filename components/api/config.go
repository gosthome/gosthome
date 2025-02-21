package api

import (
	"bytes"
	"context"
	"encoding"
	"log/slog"
	"slices"

	"github.com/gosthome/gosthome/components/api/frameshakers"
	"github.com/gosthome/gosthome/core/component"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	component.ConfigOf[Server, *Server]
	Address    string           `yaml:"address"`
	Port       uint16           `yaml:"port"`
	Password   *ConfigPassword  `yaml:"password"`
	Encryption ConfigEncryption `yaml:"encryption"`
}

func NewConfig() *Config {
	return &Config{
		Port: 6053,
	}
}

func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx, c,
		validation.Field(&c.Address),
	)
}

var _ cv.Validatable = (*Config)(nil)

type ConfigEncryption struct {
	Key *frameshakers.ConfigNoisePSK `yaml:"key"`
}

// Validate implements validation.Validatable.
func (c *ConfigEncryption) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx, c,
		validation.Field(&c.Key),
	)
}

var _ cv.Validatable = (*ConfigEncryption)(nil)

type ConfigPassword struct {
	hash []byte
}

// Validate implements validation.Validatable.
func (n *ConfigPassword) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, n, validation.Field(&n.hash, validation.Required))
}

func (n *ConfigPassword) Equal(other *ConfigPassword) bool {
	nv := n.Valid()
	ov := other.Valid()
	if nv && ov {
		return bytes.Equal(n.hash, other.hash)
	}
	return !nv && !ov
}

func (n *ConfigPassword) Valid() bool {
	if n == nil {
		return false
	}
	if len(n.hash) == 0 {
		return false
	}
	return true
}

func (n *ConfigPassword) MarshalText() ([]byte, error) {
	return slices.Clone(n.hash), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *ConfigPassword) UnmarshalText(text []byte) (err error) {
	_, err = bcrypt.Cost(text)
	if err != nil {
		n.hash, err = bcrypt.GenerateFromPassword(text, 10)
		if err != nil {
			return err
		}
		slog.Warn("Dont use plaintext password in config, please, store this password hash", "hash", string(n.hash))
	} else {
		n.hash = text
	}
	return nil
}

func ParsePassword(psk string) (*ConfigPassword, error) {
	r := &ConfigPassword{}
	err := r.UnmarshalText([]byte(psk))
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (n *ConfigPassword) Check(text string) bool {
	err := bcrypt.CompareHashAndPassword(n.hash, []byte(text))
	return err == nil
}

var _ encoding.TextUnmarshaler = (*ConfigPassword)(nil)
var _ cv.Validatable = (*ConfigPassword)(nil)
