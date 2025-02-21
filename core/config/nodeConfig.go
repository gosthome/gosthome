package config

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding"
	"fmt"
	"io"
	"net"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/goccy/go-yaml"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/registry"
)

type Config struct {
	*registry.Registry `yaml:"-"`

	Gosthome   GosthomeConfig `yaml:"gosthome"`
	Components Configs        `yaml:",inline"`
}

type lcOpt struct {
	cr *registry.Registry
}

type loadConfigOption func(*lcOpt)

func WithRegistry(reg *registry.Registry) loadConfigOption {
	return func(lo *lcOpt) {
		lo.cr = reg
	}
}

func LoadConfig(r io.Reader, opts ...loadConfigOption) (*Config, error) {
	o := lcOpt{
		cr: registry.DefaultRegistry(),
	}
	for _, opt := range opts {
		opt(&o)
	}
	ctx := context.Background()
	valid := &cv.Validator{}
	dec := yaml.NewDecoder(r, yaml.Validator(valid))
	ctx = context.WithValue(ctx, cv.ConfigYAMLDecoderKey{}, dec)
	ctx = context.WithValue(ctx, cv.ComponentRegistryKey{}, o.cr)
	valid.Context = ctx
	ret := &Config{
		Registry: o.cr,
	}
	err := dec.DecodeContext(ctx, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Validate implements validation.Validatable.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return cv.ValidateEmbedded(
		c.Components.ValidateWithContext(ctx),
		validation.ValidateStructWithContext(
			ctx, c,
			validation.Field(&c.Gosthome)),
	)
}

type GosthomeConfig struct {
	Name         string `yaml:"name"`
	FriendlyName string `yaml:"friendly_name"`
	Area         string `yaml:"area"`
	Comment      string `yaml:"comment"`
	MAC          *MAC   `yaml:"mac"`

	CompilationTime string `yaml:"compilation_time"`

	Project GosthomeProject `yaml:"project"`

	// OnBoot string
	// OnShutdown string
	// OnLoop string
}

// Validate implements validation.Validatable.
func (g GosthomeConfig) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &g,
		validation.Field(&g.Name, validation.When(g.FriendlyName == "", validation.Required)),
		validation.Field(&g.FriendlyName),
		validation.Field(&g.MAC, validation.Required),
		validation.Field(&g.Project),
	)
}

type GosthomeProject struct {
	Name    string
	Version string
}

// Validate implements validation.Validatable.
func (g *GosthomeProject) ValidateWithContext(ctx context.Context) error {
	return nil
}

type MAC struct {
	data [6]byte
}

func (n *MAC) Equal(other *MAC) bool {
	return bytes.Equal(n.data[:], other.data[:])
}

// Validate implements validation.Validatable.
func (n *MAC) ValidateWithContext(ctx context.Context) error {
	if n == nil {
		return validation.ErrNotNilRequired
	}
	if n.data == [6]byte{0, 0, 0, 0, 0, 0} {
		return validation.NewError("gosthome_zero_mac", "is all zeros")
	}
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *MAC) UnmarshalText(text []byte) error {
	macStr := string(text)
	macStr = strings.ReplaceAll(macStr, " ", "")
	macStr = strings.ToUpper(macStr)

	// Check if the MAC address is in a valid format
	mac, err := net.ParseMAC(macStr)
	if err != nil {
		return fmt.Errorf("invalid MAC address: %s", macStr)
	}
	if len(mac) > 6 {
		return fmt.Errorf("too long MAC address: %s", macStr)
	}
	copy(n.data[:], mac[:6])
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (n *MAC) MarshalText() (text []byte, err error) {
	return []byte(n.String()), nil
}

func ParseMAC(mac string) (*MAC, error) {
	r := &MAC{}
	err := r.UnmarshalText([]byte(mac))
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (n *MAC) String() string {
	return net.HardwareAddr(n.data[:]).String()
}

var _ encoding.TextMarshaler = (*MAC)(nil)
var _ encoding.TextUnmarshaler = (*MAC)(nil)

func GenerateMAC() (*MAC, error) {
	const (
		local     = 0b10
		multicast = 0b1
	)
	ret := &MAC{}
	_, err := rand.Read(ret.data[:])
	if err != nil {
		return nil, err
	}
	// clear multicast bit (&^), ensure local bit (|)
	ret.data[0] = ret.data[0]&^multicast | local
	return ret, nil
}

var _ cv.Validatable = (*Config)(nil)
var _ cv.Validatable = (*GosthomeConfig)(nil)
var _ cv.Validatable = (*GosthomeProject)(nil)
var _ cv.Validatable = (*MAC)(nil)
