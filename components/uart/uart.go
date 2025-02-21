package uart

import (
	"context"
	"log/slog"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gosthome/gosthome/core/bus"
	"github.com/gosthome/gosthome/core/component"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/entity"
	"go.bug.st/serial"
)

//go:generate go-enum

// ENUM(
// no,//disable parity control (default)
// odd,//enable odd-parity check
// even,//enable even-parity check
// mark,//enable mark-parity (always 1) check
// space,//enable space-parity (always 0) check
// )
type parity int

// ENUM(
// one, // sets 1 stop bit (default)
// one_point_five, // sets 1.5 stop bits
// two, // sets 2 stop bits
// )
type stopbits int

type Config struct {
	component.ConfigOf[UART, *UART]
	cv.MapOrValue[UARTConfig, *UARTConfig]
}

type UARTConfig struct {
	ID       string   `yaml:"id"`
	Port     string   `yaml:"port"`
	BaudRate int      `yaml:"baud_rate"`
	DataBits int      `yaml:"data_bits"`
	Parity   parity   `yaml:"parity"`
	StopBits stopbits `yaml:"stop_bits"`
}

func (c *UARTConfig) ValidateWithContext(ctx context.Context) error {
	availablePorts, err := serial.GetPortsList()
	if err != nil {
		return err
	}
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Port, cv.String(cv.OneOf(availablePorts...))),
	)
}

// Validate implements component.Config.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateWithContext(ctx, c.Configs, validation.Required, validation.Length(1, 0))
}

func NewConfig() *Config {
	return &Config{
		MapOrValue: cv.NewMapOrValue(func() *UARTConfig {
			return &UARTConfig{
				DataBits: 8,
				Parity:   ParityNo,
				StopBits: StopbitsOne,
			}
		}),
	}
}

var _ component.Config = (*Config)(nil)

type UART struct {
	idhash uint32
	port   serial.Port
	cfg    *UARTConfig
	ctx    context.Context
}

func New(ctx context.Context, cfg *Config) ([]component.Component, error) {
	ret := []component.Component{}
	b := bus.Get(ctx)
	lc := len(cfg.Configs)
	for _, uartCfg := range cfg.Configs {
		var idhash uint32
		if lc > 1 {
			id := uartCfg.ID
			if id == "" {
				id = entity.MakeID("uart")
			}
			idhash = entity.HashID(id)
		} else {
			idhash = 0
		}
		u := &UART{
			idhash: idhash,
			cfg:    uartCfg,
		}
		ret = append(ret, u)
		b.HandleServiceCalls(bus.ServiceHandler[UARTWrite](u.write))
	}
	return ret, nil
}

// Setup implements component.Component.
func (u *UART) Setup() {
	var err error
	u.port, err = serial.Open(u.cfg.Port, &serial.Mode{
		BaudRate: u.cfg.BaudRate,
		DataBits: u.cfg.DataBits,
		Parity:   serial.Parity(int(u.cfg.Parity)),
		StopBits: serial.StopBits(int(u.cfg.StopBits)),
	})
	if err != nil {
		slog.Error("Failed to initialize uart", "err", err)
		u.port = nil
		return
	}
}

type UARTRead struct {
	Key  uint32
	Data []byte
}

func (*UARTRead) EventType() string {
	return "uart_read"
}

type UARTWrite struct {
	Key  uint32
	Data []byte
}

func (*UARTWrite) ServiceType() string {
	return "uart_write"
}

func (u *UART) write(w *UARTWrite) {
	if u.idhash != w.Key {
		return
	}
	if u.port == nil {
		slog.Error("trying to write to a closed uart port")
	}
	n, err := u.port.Write(w.Data)
	if err != nil {
		slog.Error("failed to write to uart", "err", err)
	}
	slog.Debug("uart wrote n bytes", "n", n)
}

// Close implements component.Component.
func (u *UART) Close() error {
	return u.port.Close()
}

// InitializationPriority implements component.Component.
func (u *UART) InitializationPriority() component.InitializationPriority {
	return component.InitializationPriorityBus
}

var _ component.Component = (*UART)(nil)
