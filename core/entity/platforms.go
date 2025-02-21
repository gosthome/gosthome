package entity

import "github.com/gosthome/gosthome/core/component"

type BinarySensorPlatformer interface {
	BinarySensorPlatform() component.Declaration
}

type CoverPlatformer interface {
	CoverPlatform() component.Declaration
}

type FanPlatformer interface {
	FanPlatform() component.Declaration
}

type LightPlatformer interface {
	LightPlatform() component.Declaration
}

type SensorPlatformer interface {
	SensorPlatform() component.Declaration
}

type SwitchPlatformer interface {
	SwitchPlatform() component.Declaration
}

type ButtonPlatformer interface {
	ButtonPlatform() component.Declaration
}

type TextSensorPlatformer interface {
	TextSensorPlatform() component.Declaration
}

type ServicePlatformer interface {
	ServicePlatform() component.Declaration
}

type CameraPlatformer interface {
	CameraPlatform() component.Declaration
}

type ClimatePlatformer interface {
	ClimatePlatform() component.Declaration
}

type NumberPlatformer interface {
	NumberPlatform() component.Declaration
}

type DatePlatformer interface {
	DatePlatform() component.Declaration
}

type TimePlatformer interface {
	TimePlatform() component.Declaration
}

type DatetimePlatformer interface {
	DatetimePlatform() component.Declaration
}

type TextPlatformer interface {
	TextPlatform() component.Declaration
}

type SelectPlatformer interface {
	SelectPlatform() component.Declaration
}

type LockPlatformer interface {
	LockPlatform() component.Declaration
}

type ValvePlatformer interface {
	ValvePlatform() component.Declaration
}

type MediaPlayerPlatformer interface {
	MediaPlayerPlatform() component.Declaration
}

type AlarmControlPanelPlatformer interface {
	AlarmControlPanelPlatform() component.Declaration
}

type EventPlatformer interface {
	EventPlatform() component.Declaration
}

type UpdatePlatformer interface {
	UpdatePlatform() component.Declaration
}
