package entity

import (
	"context"
)

//go:generate go-enum

// ENUM(
// binary_sensor,
// cover,
// fan,
// light,
// sensor,
// switch,
// button,
// text_sensor,
// service,
// camera,
// climate,
// number,
// datetime_date,
// datetime_time,
// datetime_datetime,
// text,
// select,
// lock,
// valve,
// media_player,
// alarm_control_panel,
// event,
// update,
// )
type DomainType uint

const (
	DomainTypeStart DomainType = DomainTypeBinarySensor
	DomainTypeEnd   DomainType = DomainTypeUpdate
)

type DomainTyper interface {
	DomainType() DomainType
}

type BinarySensorDomain struct {
	BaseDomain[BinarySensorDomain, BinarySensor, *BinarySensorDomain]
}

func (*BinarySensorDomain) DomainType() DomainType {
	return DomainTypeBinarySensor
}

type BinarySensorState struct {
	State   bool
	Missing bool
}

type BinarySensor interface {
	EntityComponent
	WithState[BinarySensorState]
	WithIcon
	WithDeviceClass
}

type CoverDomain struct {
	BaseDomain[CoverDomain, Cover, *CoverDomain]
}

func (*CoverDomain) DomainType() DomainType {
	return DomainTypeCover
}

// ENUM(open,closed)
type LegacyCoverState int32

type CoverState struct {
	LegacyState LegacyCoverState
	Position    float32
	Tilt        float32
}

type Cover interface {
	EntityComponent
	WithState[CoverState]
	WithIcon
	WithDeviceClass
}

type FanDomain struct {
	BaseDomain[FanDomain, Fan, *FanDomain]
}

func (*FanDomain) DomainType() DomainType {
	return DomainTypeFan
}

// ENUM(low,medium,high)
type FanSpeed int32

// ENUM(forward.reverse)
type FanDirection int32

type FanState struct {
	State       bool
	Oscillating bool

	Speed      FanSpeed
	Direction  FanDirection
	SpeedLevel int32
	PresetMode string
}

type Fan interface {
	EntityComponent
	WithState[FanState]
	WithIcon
}

type LightDomain struct {
	BaseDomain[LightDomain, Light, *LightDomain]
}

func (*LightDomain) DomainType() DomainType {
	return DomainTypeLight
}

type LightState struct {
	State            bool
	Brightness       float32
	ColorMode        int32
	ColorBrightness  float32
	Red              float32
	Green            float32
	Blue             float32
	White            float32
	ColorTemperature float32
	ColdWhite        float32
	WarmWhite        float32
	Effect           string
}

type Light interface {
	EntityComponent
	WithState[LightState]
	WithIcon
}

type SensorDomain struct {
	BaseDomain[SensorDomain, Sensor, *SensorDomain]
}

func (*SensorDomain) DomainType() DomainType {
	return DomainTypeSensor
}

type SensorState struct {
	State        float32
	MissingState bool
}

type Sensor interface {
	EntityComponent
	WithState[SensorState]
	WithIcon
	WithDeviceClass
	WithUnitOfMeasurement
}

type SwitchDomain struct {
	BaseDomain[SwitchDomain, Switch, *SwitchDomain]
}

func (*SwitchDomain) DomainType() DomainType {
	return DomainTypeSwitch
}

type SwitchState struct {
	State bool
}

type Switch interface {
	EntityComponent
	WithState[SwitchState]
	WithIcon
	WithDeviceClass
}

type ButtonDomain struct {
	BaseDomain[ButtonDomain, Button, *ButtonDomain]
}

func (*ButtonDomain) DomainType() DomainType {
	return DomainTypeButton
}

type Button interface {
	EntityComponent
	WithIcon
	WithDeviceClass
	Press(ctx context.Context) error
}

// func NewButton() *button {
// 	ret := &button{}
// 	var _ Button = ret
// 	return ret
// }

// type button struct {
// 	BaseEntity
// 	IconMixin
// 	DeviceClassMixin
// 	onPress func(ctx context.Context, b Button) error
// }

// // Pressed implements Button.
// func (b *button) Press(ctx context.Context) error {
// 	return b.onPress(ctx, b)
// }

// func (b *button) SetOnPress(onPress func(ctx context.Context, b Button) error) error {
// 	b.onPress = onPress
// 	return nil
// }

type TextSensorDomain struct {
	BaseDomain[TextSensorDomain, TextSensor, *TextSensorDomain]
}

func (*TextSensorDomain) DomainType() DomainType {
	return DomainTypeTextSensor
}

type TextSensorState struct {
	State        string
	MissingState bool
}

type TextSensor interface {
	EntityComponent
	WithState[TextSensorState]
	WithIcon
	WithDeviceClass
}

type ServiceDomain struct {
	BaseDomain[ServiceDomain, Service, *ServiceDomain]
}

func (*ServiceDomain) DomainType() DomainType {
	return DomainTypeService
}

type Service interface {
	EntityComponent
	WithIcon
	WithDeviceClass
}

type CameraDomain struct {
	BaseDomain[CameraDomain, Camera, *CameraDomain]
}

func (*CameraDomain) DomainType() DomainType {
	return DomainTypeCamera
}

type Camera interface {
	EntityComponent
	WithIcon
}

type ClimateDomain struct {
	BaseDomain[ClimateDomain, Climate, *ClimateDomain]
}

func (*ClimateDomain) DomainType() DomainType {
	return DomainTypeClimate
}

// ENUM(off,heat_cool,cool,heat,fan_only,dry,auto)
type ClimateMode int32

// ENUM(on,off,auto,low,medium,high,middle,focus,diffuse,quiet)
type ClimateFanMode int32

// ENUM(off,both,vertical,horizontal)
type ClimateSwingMode int32

// ENUM(off,cooling,heating,idle,drying,fan)
type ClimateAction int32

// ENUM(none,home,away,boost,comfort,eco,sleep,activity)
type ClimatePreset int32

type ClimateState struct {
	Mode                  ClimateMode
	CurrentTemperature    float32
	TargetTemperature     float32
	TargetTemperatureLow  float32
	TargetTemperatureHigh float32
	// For older peers, equal to preset == CLIMATE_PRESET_AWAY
	LegacyAway      bool
	Action          ClimateAction
	FanMode         ClimateFanMode
	SwingMode       ClimateSwingMode
	CustomFanMode   string
	Preset          ClimatePreset
	CustomPreset    string
	CurrentHumidity float32
	TargetHumidity  float32
}

type Climate interface {
	EntityComponent
	WithState[ClimateState]
	WithIcon
}

type NumberDomain struct {
	BaseDomain[NumberDomain, Number, *NumberDomain]
}

func (*NumberDomain) DomainType() DomainType {
	return DomainTypeNumber
}

type NumberState struct {
	State        float32
	MissingState bool
}

type Number interface {
	EntityComponent
	WithState[NumberState]
	WithIcon
	WithDeviceClass
	WithUnitOfMeasurement
}

type DateDomain struct {
	BaseDomain[DateDomain, Date, *DateDomain]
}

func (*DateDomain) DomainType() DomainType {
	return DomainTypeDatetimeDate
}

type DateState struct {
	MissingState bool
	Year         uint32
	Month        uint32
	Day          uint32
}

type Date interface {
	EntityComponent
	WithState[DateState]
	WithIcon
}

type TimeDomain struct {
	BaseDomain[TimeDomain, Time, *TimeDomain]
}

func (*TimeDomain) DomainType() DomainType {
	return DomainTypeDatetimeTime
}

type TimeState struct {
	MissingState bool
	Hour         uint32
	Minute       uint32
	Second       uint32
}

type Time interface {
	EntityComponent
	WithState[TimeState]
	WithIcon
}

type DatetimeDomain struct {
	BaseDomain[DatetimeDomain, Datetime, *DatetimeDomain]
}

func (*DatetimeDomain) DomainType() DomainType {
	return DomainTypeDatetimeDatetime
}

type DatetimeState struct {
	MissingState bool
	EpochSeconds uint32
}

type Datetime interface {
	EntityComponent
	WithState[DatetimeState]
	WithIcon
}

type TextDomain struct {
	BaseDomain[TextDomain, Text, *TextDomain]
}

func (*TextDomain) DomainType() DomainType {
	return DomainTypeText
}

type TextState struct {
	State        string
	MissingState bool
}

type Text interface {
	EntityComponent
	WithState[TextState]
	WithIcon
}

type SelectDomain struct {
	BaseDomain[SelectDomain, Select, *SelectDomain]
}

func (*SelectDomain) DomainType() DomainType {
	return DomainTypeSelect
}

type SelectState struct {
	State        string
	MissingState bool
}

type Select interface {
	EntityComponent
	WithState[SelectState]
	WithIcon
}

type LockDomain struct {
	BaseDomain[LockDomain, Lock, *LockDomain]
}

func (*LockDomain) DomainType() DomainType {
	return DomainTypeLock
}

// ENUM(none,locked,unlocked,jammed,locking,unlocking)
type LockState int32

type Lock interface {
	EntityComponent
	WithState[LockState]
	WithIcon
}

type ValveDomain struct {
	BaseDomain[ValveDomain, Valve, *ValveDomain]
}

func (*ValveDomain) DomainType() DomainType {
	return DomainTypeValve
}

// ENUM(idle,is_opening,is_closing)
type ValveOperation int32

type ValveState struct {
	Position         float32
	CurrentOperation ValveOperation
}

type Valve interface {
	EntityComponent
	WithState[ValveState]
	WithIcon
	WithDeviceClass
}

type MediaPlayerDomain struct {
	BaseDomain[MediaPlayerDomain, MediaPlayer, *MediaPlayerDomain]
}

func (*MediaPlayerDomain) DomainType() DomainType {
	return DomainTypeMediaPlayer
}

// ENUM(none,idle,playing,paused)
type MediaPlayingState int32

type MediaPlayerState struct {
	State  MediaPlayingState
	Volume float32
	Muted  bool
}

type MediaPlayer interface {
	EntityComponent
	WithState[MediaPlayerState]
	WithIcon
}

type AlarmControlPanelDomain struct {
	BaseDomain[AlarmControlPanelDomain, AlarmControlPanel, *AlarmControlPanelDomain]
}

func (*AlarmControlPanelDomain) DomainType() DomainType {
	return DomainTypeAlarmControlPanel
}

// ENUM(disarmed,armed_home,armed_away,armed_night,armed_vacation,armed_custom_bypass,pending,arming,disarming,triggered)
type AlarmControlPanelState int32

type AlarmControlPanel interface {
	EntityComponent
	WithState[AlarmControlPanelState]
	WithIcon
}

type EventDomain struct {
	BaseDomain[EventDomain, Event, *EventDomain]
}

func (*EventDomain) DomainType() DomainType {
	return DomainTypeEvent
}

type Event interface {
	EntityComponent
	WithIcon
	WithDeviceClass
}

type UpdateDomain struct {
	BaseDomain[UpdateDomain, Update, *UpdateDomain]
}

func (*UpdateDomain) DomainType() DomainType {
	return DomainTypeUpdate
}

type UpdateState struct {
	MissingState   bool
	InProgress     bool
	HasProgress    bool
	Progress       float32
	CurrentVersion string
	LatestVersion  string
	Title          string
	ReleaseSummary string
	ReleaseUrl     string
}

type Update interface {
	EntityComponent
	WithState[UpdateState]
	WithIcon
	WithDeviceClass
}
