package registry

import (
	"errors"

	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/entity"
)

func (r *Registry) tryRegisterPlatforms(name string, cd component.Declaration) error {
	errs := []error{}
	if p, ok := cd.(entity.BinarySensorPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeBinarySensor, name, p.BinarySensorPlatform()))
	}
	if p, ok := cd.(entity.CoverPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeCover, name, p.CoverPlatform()))
	}
	if p, ok := cd.(entity.FanPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeFan, name, p.FanPlatform()))
	}
	if p, ok := cd.(entity.LightPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeLight, name, p.LightPlatform()))
	}
	if p, ok := cd.(entity.SensorPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeSensor, name, p.SensorPlatform()))
	}
	if p, ok := cd.(entity.SwitchPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeSwitch, name, p.SwitchPlatform()))
	}
	if p, ok := cd.(entity.ButtonPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeButton, name, p.ButtonPlatform()))
	}
	if p, ok := cd.(entity.TextSensorPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeTextSensor, name, p.TextSensorPlatform()))
	}
	if p, ok := cd.(entity.ServicePlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeService, name, p.ServicePlatform()))
	}
	if p, ok := cd.(entity.CameraPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeCamera, name, p.CameraPlatform()))
	}
	if p, ok := cd.(entity.ClimatePlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeClimate, name, p.ClimatePlatform()))
	}
	if p, ok := cd.(entity.NumberPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeNumber, name, p.NumberPlatform()))
	}
	if p, ok := cd.(entity.DatePlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeDatetimeDate, name, p.DatePlatform()))
	}
	if p, ok := cd.(entity.TimePlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeDatetimeTime, name, p.TimePlatform()))
	}
	if p, ok := cd.(entity.DatetimePlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeDatetimeDatetime, name, p.DatetimePlatform()))
	}
	if p, ok := cd.(entity.TextPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeText, name, p.TextPlatform()))
	}
	if p, ok := cd.(entity.SelectPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeSelect, name, p.SelectPlatform()))
	}
	if p, ok := cd.(entity.LockPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeLock, name, p.LockPlatform()))
	}
	if p, ok := cd.(entity.ValvePlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeValve, name, p.ValvePlatform()))
	}
	if p, ok := cd.(entity.MediaPlayerPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeMediaPlayer, name, p.MediaPlayerPlatform()))
	}
	if p, ok := cd.(entity.AlarmControlPanelPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeAlarmControlPanel, name, p.AlarmControlPanelPlatform()))
	}
	if p, ok := cd.(entity.EventPlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeEvent, name, p.EventPlatform()))
	}
	if p, ok := cd.(entity.UpdatePlatformer); ok {
		errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeUpdate, name, p.UpdatePlatform()))
	}
	return errors.Join(errs...)
}
