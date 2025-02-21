package config_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
	"unsafe"

	"github.com/gkampitakis/go-snaps/snaps"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/goccy/go-yaml"
	"github.com/google/go-cmp/cmp"
	"github.com/gosthome/gosthome/components"
	"github.com/gosthome/gosthome/components/api"
	"github.com/gosthome/gosthome/components/api/frameshakers"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/config"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/registry"
	"github.com/matryer/is"
)

func file(t *testing.T, fn string) io.Reader {
	t.Helper()
	f, err := os.Open(fn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { f.Close() })
	return f
}

type testComponent struct{}

// Setup implements component.Component.
func (t *testComponent) Setup() {
}

// InitializationPriority implements component.Component.
func (t *testComponent) InitializationPriority() component.InitializationPriority {
	return component.InitializationPriorityBus
}

func (t *testComponent) Close() error {
	return nil
}

var _ component.Component = (*testComponent)(nil)

type testComponentConfig struct {
	component.ConfigOf[testComponent, *testComponent]
	A string `yaml:"a"`
}

// Validate implements component.ComponentConfig.
func (t *testComponentConfig) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, t, validation.Field(&t.A, validation.Required, validation.By(func(value interface{}) error {
		if value.(string) == "wrong value" {
			return validation.NewError("test_wrong_value", "a is wrong")
		}
		return nil
	})))
}

var _ component.Config = (*testComponentConfig)(nil)

type testComponentDeclaration struct{}

// Component implements component.ComponentDeclaration.
func (t *testComponentDeclaration) Component(ctx context.Context, conf component.Config) ([]component.Component, error) {
	_ = conf.(*testComponentConfig)
	return []component.Component{&testComponent{}}, nil
}

// Config implements component.ComponentDeclaration.
func (t *testComponentDeclaration) Config() *component.ConfigDecoder {
	return component.NewConfigDecoder(&testComponentConfig{})
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func TestYaml(t *testing.T) {
	is := is.New(t)

	type testdata struct {
		name     string
		data     io.Reader
		isErr    bool
		expected *testComponentConfig
	}

	tests := []testdata{
		{
			name:  "ok",
			isErr: false,
			data: bytes.NewBuffer([]byte(`
a: b
`)),
			expected: &testComponentConfig{A: "b"},
		},
		{
			name:  "notok",
			isErr: true,
			data: bytes.NewBuffer([]byte(`
a: "wrong value"
`)),
			expected: &testComponentConfig{A: "b"},
		},
	}
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			is := is.New(t)
			actual := &testComponentConfig{}
			dec := yaml.NewDecoder(tcase.data, yaml.Validator(yaml.StructValidator(&cv.Validator{context.Background()})))
			err := dec.Decode(&actual)
			if tcase.isErr {
				is.True(err != nil)
				snaps.MatchSnapshot(t, err.Error())
			} else {
				if diff := cmp.Diff(actual, tcase.expected); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	is := is.New(t)

	type testdata struct {
		name  string
		data  io.Reader
		isErr bool

		reg *registry.Registry

		expected *config.Config
	}

	tests := []testdata{
		{
			name:  "gibberish",
			data:  bytes.NewBuffer([]byte("jjjhj")),
			isErr: true,
		},
		{
			name: "only-gosthome",
			data: bytes.NewBuffer([]byte(`
gosthome:
`)),
			isErr: true,
		},
		{
			name: "valid_gosthome_with_unregistered_component",
			data: bytes.NewBuffer([]byte(`
gosthome:
  name: exampl
  mac: 00:aa:bb:cc:dd:ee

app:
  a: b
`)),
			isErr: true,
		},
		{
			name: "valid_gosthome_with_unregistered_empty_component",
			data: bytes.NewBuffer([]byte(`
gosthome:
  name: exampl
  mac: 00:aa:bb:cc:dd:ee

app:
`)),
			isErr: true,
		},
		{
			name: "valid_gosthome_with_registered_component_and_invalid_component",
			data: bytes.NewBuffer([]byte(`
gosthome:
  name: exampl
  mac: 00:aa:bb:cc:dd:ee

app:
  a: "wrong value"
`)),
			isErr: true,
			reg: func() *registry.Registry {
				rg := registry.NewRegistry()
				rg.Register("app", &testComponentDeclaration{})
				return rg
			}(),
		},
		{
			name: "valid_gosthome_with_registered_component",
			data: bytes.NewBuffer([]byte(`
gosthome:
  name: exampl
  mac: 00:aa:bb:cc:dd:ee

app:
  a: b
`)),
			isErr: false,
			reg: func() *registry.Registry {
				rg := registry.NewRegistry()
				rg.Register("app", &testComponentDeclaration{})
				return rg
			}(),
			expected: &config.Config{
				Registry: &registry.Registry{},
				Gosthome: config.GosthomeConfig{Name: "exampl", MAC: must(config.ParseMAC("00:aa:bb:cc:dd:ee"))},
				Components: config.Configs{"app": func() *component.ConfigDecoder {
					return &component.ConfigDecoder{Config: &testComponentConfig{A: "b"}}
				}()},
			},
		},
		{
			name: "valid_gosthome_with_api",
			data: bytes.NewBuffer([]byte(`
gosthome:
  name: exampl
  mac: 00:aa:bb:cc:dd:ee

api:
  address: 127.0.0.1
  port: 6969
  encryption:
    key: 9kD0vcdCbh9UQWaSCUJXsX3Rt0PWj5BHWoqMTI2TTkM=
`)),
			isErr: false,
			reg: func() *registry.Registry {
				rg := registry.NewRegistry()
				components.RegisterAll(rg)
				return rg
			}(),
			expected: &config.Config{
				Registry: &registry.Registry{},
				Gosthome: config.GosthomeConfig{Name: "exampl", MAC: must(config.ParseMAC("00:aa:bb:cc:dd:ee"))},
				Components: config.Configs{"api": func() *component.ConfigDecoder {
					return &component.ConfigDecoder{Config: &api.Config{
						Address: "127.0.0.1",
						Port:    6969,
						Encryption: api.ConfigEncryption{
							Key: must(frameshakers.ParseNoisePSK("9kD0vcdCbh9UQWaSCUJXsX3Rt0PWj5BHWoqMTI2TTkM=")),
						},
					}}
				}()},
			},
		},
	}
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			is := is.New(t)
			if tcase.reg == nil {
				tcase.reg = &registry.Registry{}
			}
			c, err := config.LoadConfig(tcase.data, config.WithRegistry(tcase.reg))
			if tcase.isErr {
				is.Equal(c, (*config.Config)(nil))
				is.True(err != nil)
				snaps.MatchSnapshot(t, err.Error())
			} else {
				is.NoErr(err)
				tcase.expected.Registry = tcase.reg
				if diff := cmp.Diff(
					c, tcase.expected,
					cmp.Comparer(func(rl, rr *registry.Registry) bool {
						return uintptr(unsafe.Pointer(rl)) == uintptr(unsafe.Pointer(rr))
					}),
					cmp.FilterPath(
						func(p cmp.Path) bool {
							if len(p) != 6 {
								return false
							}
							sf, ok := p.Last().(cmp.StructField)
							if !ok {
								return false
							}
							return sf.Name() == "Unmarshal" || sf.Name() == "Marshal"
						},
						cmp.Ignore(),
					)); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
