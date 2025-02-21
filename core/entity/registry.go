package entity

import (
	"errors"
	"log/slog"
)

type Registry struct {
	BinarySensorDomain      *BinarySensorDomain
	CoverDomain             *CoverDomain
	FanDomain               *FanDomain
	LightDomain             *LightDomain
	SensorDomain            *SensorDomain
	SwitcheDomain           *SwitchDomain
	ButtonDomain            *ButtonDomain
	TextSensorDomain        *TextSensorDomain
	ServiceDomain           *ServiceDomain
	CameraDomain            *CameraDomain
	ClimateDomain           *ClimateDomain
	NumberDomain            *NumberDomain
	DateDomain              *DateDomain
	TimeDomain              *TimeDomain
	DatetimeDomain          *DatetimeDomain
	TextDomain              *TextDomain
	SelectDomain            *SelectDomain
	LockDomain              *LockDomain
	ValveDomain             *ValveDomain
	MediaPlayerDomain       *MediaPlayerDomain
	AlarmControlPanelDomain *AlarmControlPanelDomain
	EventDomain             *EventDomain
	UpdateDomain            *UpdateDomain
}

func (er *Registry) BinarySensorByKey(key uint32) (BinarySensor, bool) {
	if er.BinarySensorDomain == nil {
		slog.Error("BinarySensorDomain is not registered!")
		return nil, false
	}
	return er.BinarySensorDomain.FindByKey(key)
}
func (er *Registry) BinarySensors() []Entity {
	if er.BinarySensorDomain == nil {
		slog.Debug("BinarySensorDomain is not registered!")
		return nil
	}
	return er.BinarySensorDomain.Clone()
}
func (er *Registry) CoverByKey(key uint32) (Cover, bool) {
	if er.CoverDomain == nil {
		slog.Error("CoverDomain is not registered!")
		return nil, false
	}
	return er.CoverDomain.FindByKey(key)
}
func (er *Registry) Covers() []Entity {
	if er.CoverDomain == nil {
		slog.Debug("CoverDomain is not registered!")
		return nil
	}
	return er.CoverDomain.Clone()
}
func (er *Registry) FanByKey(key uint32) (Fan, bool) {
	if er.FanDomain == nil {
		slog.Error("FanDomain is not registered!")
		return nil, false
	}
	return er.FanDomain.FindByKey(key)
}
func (er *Registry) Fans() []Entity {
	if er.FanDomain == nil {
		slog.Debug("FanDomain is not registered!")
		return nil
	}
	return er.FanDomain.Clone()
}
func (er *Registry) LightByKey(key uint32) (Light, bool) {
	if er.LightDomain == nil {
		slog.Error("LightDomain is not registered!")
		return nil, false
	}
	return er.LightDomain.FindByKey(key)
}
func (er *Registry) Lights() []Entity {
	if er.LightDomain == nil {
		slog.Debug("LightDomain is not registered!")
		return nil
	}
	return er.LightDomain.Clone()
}
func (er *Registry) SensorByKey(key uint32) (Sensor, bool) {
	if er.SensorDomain == nil {
		slog.Error("SensorDomain is not registered!")
		return nil, false
	}
	return er.SensorDomain.FindByKey(key)
}
func (er *Registry) Sensors() []Entity {
	if er.SensorDomain == nil {
		slog.Debug("SensorDomain is not registered!")
		return nil
	}
	return er.SensorDomain.Clone()
}
func (er *Registry) SwitchByKey(key uint32) (Switch, bool) {
	if er.SwitcheDomain == nil {
		slog.Error("SwitcheDomain is not registered!")
		return nil, false
	}
	return er.SwitcheDomain.FindByKey(key)
}
func (er *Registry) Switches() []Entity {
	if er.SwitcheDomain == nil {
		slog.Debug("SwitcheDomain is not registered!")
		return nil
	}
	return er.SwitcheDomain.Clone()
}
func (er *Registry) ButtonByKey(key uint32) (Button, bool) {
	if er.ButtonDomain == nil {
		slog.Error("ButtonDomain is not registered!")
		return nil, false
	}
	return er.ButtonDomain.FindByKey(key)
}
func (er *Registry) Buttons() []Entity {
	if er.ButtonDomain == nil {
		slog.Debug("ButtonDomain is not registered!")
		return nil
	}
	return er.ButtonDomain.Clone()
}
func (er *Registry) TextSensorByKey(key uint32) (TextSensor, bool) {
	if er.TextSensorDomain == nil {
		slog.Error("TextSensorDomain is not registered!")
		return nil, false
	}
	return er.TextSensorDomain.FindByKey(key)
}
func (er *Registry) TextSensors() []Entity {
	if er.TextSensorDomain == nil {
		slog.Debug("TextSensorDomain is not registered!")
		return nil
	}
	return er.TextSensorDomain.Clone()
}
func (er *Registry) ServiceByKey(key uint32) (Service, bool) {
	if er.ServiceDomain == nil {
		slog.Error("ServiceDomain is not registered!")
		return nil, false
	}
	return er.ServiceDomain.FindByKey(key)
}
func (er *Registry) Services() []Entity {
	if er.ServiceDomain == nil {
		slog.Debug("ServiceDomain is not registered!")
		return nil
	}
	return er.ServiceDomain.Clone()
}
func (er *Registry) CameraByKey(key uint32) (Camera, bool) {
	if er.CameraDomain == nil {
		slog.Error("CameraDomain is not registered!")
		return nil, false
	}
	return er.CameraDomain.FindByKey(key)
}
func (er *Registry) Cameras() []Entity {
	if er.CameraDomain == nil {
		slog.Debug("CameraDomain is not registered!")
		return nil
	}
	return er.CameraDomain.Clone()
}
func (er *Registry) ClimateByKey(key uint32) (Climate, bool) {
	if er.ClimateDomain == nil {
		slog.Error("ClimateDomain is not registered!")
		return nil, false
	}
	return er.ClimateDomain.FindByKey(key)
}
func (er *Registry) Climates() []Entity {
	if er.ClimateDomain == nil {
		slog.Debug("ClimateDomain is not registered!")
		return nil
	}
	return er.ClimateDomain.Clone()
}
func (er *Registry) NumberByKey(key uint32) (Number, bool) {
	if er.NumberDomain == nil {
		slog.Error("NumberDomain is not registered!")
		return nil, false
	}
	return er.NumberDomain.FindByKey(key)
}
func (er *Registry) Numbers() []Entity {
	if er.NumberDomain == nil {
		slog.Debug("NumberDomain is not registered!")
		return nil
	}
	return er.NumberDomain.Clone()
}
func (er *Registry) DateByKey(key uint32) (Date, bool) {
	if er.DateDomain == nil {
		slog.Error("DateDomain is not registered!")
		return nil, false
	}
	return er.DateDomain.FindByKey(key)
}
func (er *Registry) Dates() []Entity {
	if er.DateDomain == nil {
		slog.Debug("DateDomain is not registered!")
		return nil
	}
	return er.DateDomain.Clone()
}
func (er *Registry) TimeByKey(key uint32) (Time, bool) {
	if er.TimeDomain == nil {
		slog.Error("TimeDomain is not registered!")
		return nil, false
	}
	return er.TimeDomain.FindByKey(key)
}
func (er *Registry) Times() []Entity {
	if er.TimeDomain == nil {
		slog.Debug("TimeDomain is not registered!")
		return nil
	}
	return er.TimeDomain.Clone()
}
func (er *Registry) DatetimeByKey(key uint32) (Datetime, bool) {
	if er.DatetimeDomain == nil {
		slog.Error("DatetimeDomain is not registered!")
		return nil, false
	}
	return er.DatetimeDomain.FindByKey(key)
}
func (er *Registry) Datetimes() []Entity {
	if er.DatetimeDomain == nil {
		slog.Debug("DatetimeDomain is not registered!")
		return nil
	}
	return er.DatetimeDomain.Clone()
}
func (er *Registry) TextByKey(key uint32) (Text, bool) {
	if er.TextDomain == nil {
		slog.Error("TextDomain is not registered!")
		return nil, false
	}
	return er.TextDomain.FindByKey(key)
}
func (er *Registry) Texts() []Entity {
	if er.TextDomain == nil {
		slog.Debug("TextDomain is not registered!")
		return nil
	}
	return er.TextDomain.Clone()
}
func (er *Registry) SelectByKey(key uint32) (Select, bool) {
	if er.SelectDomain == nil {
		slog.Error("SelectDomain is not registered!")
		return nil, false
	}
	return er.SelectDomain.FindByKey(key)
}
func (er *Registry) Selects() []Entity {
	if er.SelectDomain == nil {
		slog.Debug("SelectDomain is not registered!")
		return nil
	}
	return er.SelectDomain.Clone()
}
func (er *Registry) LockByKey(key uint32) (Lock, bool) {
	if er.LockDomain == nil {
		slog.Error("LockDomain is not registered!")
		return nil, false
	}
	return er.LockDomain.FindByKey(key)
}
func (er *Registry) Locks() []Entity {
	if er.LockDomain == nil {
		slog.Debug("LockDomain is not registered!")
		return nil
	}
	return er.LockDomain.Clone()
}
func (er *Registry) ValveByKey(key uint32) (Valve, bool) {
	if er.ValveDomain == nil {
		slog.Error("ValveDomain is not registered!")
		return nil, false
	}
	return er.ValveDomain.FindByKey(key)
}
func (er *Registry) Valves() []Entity {
	if er.ValveDomain == nil {
		slog.Debug("ValveDomain is not registered!")
		return nil
	}
	return er.ValveDomain.Clone()
}
func (er *Registry) MediaPlayerByKey(key uint32) (MediaPlayer, bool) {
	if er.MediaPlayerDomain == nil {
		slog.Error("MediaPlayerDomain is not registered!")
		return nil, false
	}
	return er.MediaPlayerDomain.FindByKey(key)
}
func (er *Registry) MediaPlayers() []Entity {
	if er.MediaPlayerDomain == nil {
		slog.Debug("MediaPlayerDomain is not registered!")
		return nil
	}
	return er.MediaPlayerDomain.Clone()
}
func (er *Registry) AlarmControlPanelByKey(key uint32) (AlarmControlPanel, bool) {
	if er.AlarmControlPanelDomain == nil {
		slog.Error("AlarmControlPanelDomain is not registered!")
		return nil, false
	}
	return er.AlarmControlPanelDomain.FindByKey(key)
}
func (er *Registry) AlarmControlPanels() []Entity {
	if er.AlarmControlPanelDomain == nil {
		slog.Debug("AlarmControlPanelDomain is not registered!")
		return nil
	}
	return er.AlarmControlPanelDomain.Clone()
}
func (er *Registry) EventByKey(key uint32) (Event, bool) {
	if er.EventDomain == nil {
		slog.Error("EventDomain is not registered!")
		return nil, false
	}
	return er.EventDomain.FindByKey(key)
}
func (er *Registry) Events() []Entity {
	if er.EventDomain == nil {
		slog.Debug("EventDomain is not registered!")
		return nil
	}
	return er.EventDomain.Clone()
}
func (er *Registry) UpdateByKey(key uint32) (Update, bool) {
	if er.UpdateDomain == nil {
		slog.Error("UpdateDomain is not registered!")
		return nil, false
	}
	return er.UpdateDomain.FindByKey(key)
}
func (er *Registry) Updates() []Entity {
	if er.UpdateDomain == nil {
		slog.Debug("UpdateDomain is not registered!")
		return nil
	}
	return er.UpdateDomain.Clone()
}

func (er *Registry) RegisterBinarySensor(ent BinarySensor) (err error) {
	if er.BinarySensorDomain == nil {
		return errors.New("BinarySensorDomain is not registered!")
	}
	return er.BinarySensorDomain.Register(ent)
}

func (er *Registry) RegisterCover(ent Cover) (err error) {
	if er.CoverDomain == nil {
		return errors.New("CoverDomain is not registered!")
	}
	return er.CoverDomain.Register(ent)
}

func (er *Registry) RegisterFan(ent Fan) (err error) {
	if er.FanDomain == nil {
		return errors.New("FanDomain is not registered!")
	}
	return er.FanDomain.Register(ent)
}

func (er *Registry) RegisterLight(ent Light) (err error) {
	if er.LightDomain == nil {
		return errors.New("LightDomain is not registered!")
	}
	return er.LightDomain.Register(ent)
}

func (er *Registry) RegisterSensor(ent Sensor) (err error) {
	if er.SensorDomain == nil {
		return errors.New("SensorDomain is not registered!")
	}
	return er.SensorDomain.Register(ent)
}

func (er *Registry) RegisterSwitch(ent Switch) (err error) {
	if er.SwitcheDomain == nil {
		return errors.New("SwitcheDomain is not registered!")
	}
	return er.SwitcheDomain.Register(ent)
}

func (er *Registry) RegisterButton(ent Button) (err error) {
	if er.ButtonDomain == nil {
		return errors.New("ButtonDomain is not registered!")
	}
	return er.ButtonDomain.Register(ent)
}

func (er *Registry) RegisterTextSensor(ent TextSensor) (err error) {
	if er.TextSensorDomain == nil {
		return errors.New("TextSensorDomain is not registered!")
	}
	return er.TextSensorDomain.Register(ent)
}

func (er *Registry) RegisterService(ent Camera) (err error) {
	if er.CameraDomain == nil {
		return errors.New("CameraDomain is not registered!")
	}
	return er.CameraDomain.Register(ent)
}

func (er *Registry) RegisterCamera(ent Camera) (err error) {
	if er.CameraDomain == nil {
		return errors.New("CameraDomain is not registered!")
	}
	return er.CameraDomain.Register(ent)
}

func (er *Registry) RegisterClimate(ent Climate) (err error) {
	if er.ClimateDomain == nil {
		return errors.New("ClimateDomain is not registered!")
	}
	return er.ClimateDomain.Register(ent)
}

func (er *Registry) RegisterNumber(ent Number) (err error) {
	if er.NumberDomain == nil {
		return errors.New("NumberDomain is not registered!")
	}
	return er.NumberDomain.Register(ent)
}

func (er *Registry) RegisterDate(ent Date) (err error) {
	if er.DateDomain == nil {
		return errors.New("DateDomain is not registered!")
	}
	return er.DateDomain.Register(ent)
}

func (er *Registry) RegisterTime(ent Time) (err error) {
	if er.TimeDomain == nil {
		return errors.New("TimeDomain is not registered!")
	}
	return er.TimeDomain.Register(ent)
}

func (er *Registry) RegisterDatetime(ent Datetime) (err error) {
	if er.DatetimeDomain == nil {
		return errors.New("DatetimeDomain is not registered!")
	}
	return er.DatetimeDomain.Register(ent)
}

func (er *Registry) RegisterText(ent Text) (err error) {
	if er.TextDomain == nil {
		return errors.New("TextDomain is not registered!")
	}
	return er.TextDomain.Register(ent)
}

func (er *Registry) RegisterSelect(ent Select) (err error) {
	if er.SelectDomain == nil {
		return errors.New("SelectDomain is not registered!")
	}
	return er.SelectDomain.Register(ent)
}

func (er *Registry) RegisterLock(ent Lock) (err error) {
	if er.LockDomain == nil {
		return errors.New("LockDomain is not registered!")
	}
	return er.LockDomain.Register(ent)
}

func (er *Registry) RegisterValve(ent Valve) (err error) {
	if er.ValveDomain == nil {
		return errors.New("ValveDomain is not registered!")
	}
	return er.ValveDomain.Register(ent)
}

func (er *Registry) RegisterMediaPlayer(ent MediaPlayer) (err error) {
	if er.MediaPlayerDomain == nil {
		return errors.New("MediaPlayerDomain is not registered!")
	}
	return er.MediaPlayerDomain.Register(ent)
}

func (er *Registry) RegisterAlarmControlPanel(ent AlarmControlPanel) (err error) {
	if er.AlarmControlPanelDomain == nil {
		return errors.New("AlarmControlPanelDomain is not registered!")
	}
	return er.AlarmControlPanelDomain.Register(ent)
}

func (er *Registry) RegisterEvent(ent Event) (err error) {
	if er.EventDomain == nil {
		return errors.New("EventDomain is not registered!")
	}
	return er.EventDomain.Register(ent)
}

func (er *Registry) RegisterUpdate(ent Update) (err error) {
	if er.UpdateDomain == nil {
		return errors.New("UpdateDomain is not registered!")
	}
	return er.UpdateDomain.Register(ent)
}
