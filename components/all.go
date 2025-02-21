// this file is autogenerated by gen.sh DO NOT EDIT.
package components

//go:generate ./gen.sh

import (
	"context"

	"github.com/gosthome/gosthome/components/api"
	"github.com/gosthome/gosthome/components/binarysensor"
	"github.com/gosthome/gosthome/components/button"
	"github.com/gosthome/gosthome/components/demo"
	"github.com/gosthome/gosthome/components/file"
	"github.com/gosthome/gosthome/components/psutil"
	"github.com/gosthome/gosthome/components/sensor"
	"github.com/gosthome/gosthome/components/textsensor"
	"github.com/gosthome/gosthome/components/uart"
	"github.com/gosthome/gosthome/components/webserver"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/registry"
)

type apiComponent struct{}

func (apiComponent) Config() *component.ConfigDecoder {
	return component.NewConfigDecoder(api.NewConfig())
}

func (apiComponent) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {
	apiCfg := cfg.(*api.Config)
	return api.New(ctx, apiCfg)
}

type binarysensorComponent struct{}

func (binarysensorComponent) Config() *component.ConfigDecoder {
	return component.NewConfigDecoder(binarysensor.NewConfig())
}

func (binarysensorComponent) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {
	binarysensorCfg := cfg.(*binarysensor.Config)
	return binarysensor.New(ctx, binarysensorCfg)
}

type buttonComponent struct{}

func (buttonComponent) Config() *component.ConfigDecoder {
	return component.NewConfigDecoder(button.NewConfig())
}

func (buttonComponent) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {
	buttonCfg := cfg.(*button.Config)
	return button.New(ctx, buttonCfg)
}

type demoComponent struct{}

func (demoComponent) Config() *component.ConfigDecoder {
	return component.NewConfigDecoder(demo.NewConfig())
}

func (demoComponent) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {
	demoCfg := cfg.(*demo.Config)
	return demo.New(ctx, demoCfg)
}

type fileComponent struct{}

func (fileComponent) Config() *component.ConfigDecoder {
	return component.NewConfigDecoder(file.NewConfig())
}

func (fileComponent) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {
	fileCfg := cfg.(*file.Config)
	return file.New(ctx, fileCfg)
}

type psutilComponent struct{}

func (psutilComponent) Config() *component.ConfigDecoder {
	return component.NewConfigDecoder(psutil.NewConfig())
}

func (psutilComponent) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {
	psutilCfg := cfg.(*psutil.Config)
	return psutil.New(ctx, psutilCfg)
}

type sensorComponent struct{}

func (sensorComponent) Config() *component.ConfigDecoder {
	return component.NewConfigDecoder(sensor.NewConfig())
}

func (sensorComponent) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {
	sensorCfg := cfg.(*sensor.Config)
	return sensor.New(ctx, sensorCfg)
}

type textsensorComponent struct{}

func (textsensorComponent) Config() *component.ConfigDecoder {
	return component.NewConfigDecoder(textsensor.NewConfig())
}

func (textsensorComponent) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {
	textsensorCfg := cfg.(*textsensor.Config)
	return textsensor.New(ctx, textsensorCfg)
}

type uartComponent struct{}

func (uartComponent) Config() *component.ConfigDecoder {
	return component.NewConfigDecoder(uart.NewConfig())
}

func (uartComponent) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {
	uartCfg := cfg.(*uart.Config)
	return uart.New(ctx, uartCfg)
}

type uartButtonEntityComponent struct{}

func (uartButtonEntityComponent) Config() *component.ConfigDecoder {
	return component.NewConfigDecoder(uart.NewButtonConfig())
}

func (uartButtonEntityComponent) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {
	uartButtonCfg := cfg.(*uart.ButtonConfig)
	return uart.NewButton(ctx, uartButtonCfg)
}

func (uartComponent) ButtonPlatform() component.Declaration {
	return &uartButtonEntityComponent{}
}

type webserverComponent struct{}

func (webserverComponent) Config() *component.ConfigDecoder {
	return component.NewConfigDecoder(webserver.NewConfig())
}

func (webserverComponent) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {
	webserverCfg := cfg.(*webserver.Config)
	return webserver.New(ctx, webserverCfg)
}

var (
	COMPONENT_KEY_API          = "api"
	COMPONENT_KEY_BINARYSENSOR = binarysensor.COMPONENT_KEY
	COMPONENT_KEY_BUTTON       = button.COMPONENT_KEY
	COMPONENT_KEY_DEMO         = "demo"
	COMPONENT_KEY_FILE         = "file"
	COMPONENT_KEY_PSUTIL       = "psutil"
	COMPONENT_KEY_SENSOR       = sensor.COMPONENT_KEY
	COMPONENT_KEY_TEXTSENSOR   = textsensor.COMPONENT_KEY
	COMPONENT_KEY_UART         = uart.COMPONENT_KEY
	COMPONENT_KEY_WEBSERVER    = webserver.COMPONENT_KEY
)

var (
	_ = registry.RegisterDefaultComponent(COMPONENT_KEY_API, apiComponent{})
	_ = registry.RegisterDefaultComponent(COMPONENT_KEY_BINARYSENSOR, binarysensorComponent{})
	_ = registry.RegisterDefaultComponent(COMPONENT_KEY_BUTTON, buttonComponent{})
	_ = registry.RegisterDefaultComponent(COMPONENT_KEY_DEMO, demoComponent{})
	_ = registry.RegisterDefaultComponent(COMPONENT_KEY_FILE, fileComponent{})
	_ = registry.RegisterDefaultComponent(COMPONENT_KEY_PSUTIL, psutilComponent{})
	_ = registry.RegisterDefaultComponent(COMPONENT_KEY_SENSOR, sensorComponent{})
	_ = registry.RegisterDefaultComponent(COMPONENT_KEY_TEXTSENSOR, textsensorComponent{})
	_ = registry.RegisterDefaultComponent(COMPONENT_KEY_UART, uartComponent{})
	_ = registry.RegisterDefaultComponent(COMPONENT_KEY_WEBSERVER, webserverComponent{})
)
