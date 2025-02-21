package entity

import "iter"

func iteratorTypes(reg *Registry) iter.Seq2[DomainType, func() []Entity] {
	return func(yield func(DomainType, func() []Entity) bool) {
		if !yield(DomainTypeBinarySensor, reg.BinarySensors) {
			return
		}
		if !yield(DomainTypeCover, reg.Covers) {
			return
		}
		if !yield(DomainTypeFan, reg.Fans) {
			return
		}
		if !yield(DomainTypeLight, reg.Lights) {
			return
		}
		if !yield(DomainTypeSensor, reg.Sensors) {
			return
		}
		if !yield(DomainTypeSwitch, reg.Switches) {
			return
		}
		if !yield(DomainTypeButton, reg.Buttons) {
			return
		}
		if !yield(DomainTypeTextSensor, reg.TextSensors) {
			return
		}
		if !yield(DomainTypeService, reg.Services) {
			return
		}
		if !yield(DomainTypeCamera, reg.Cameras) {
			return
		}
		if !yield(DomainTypeClimate, reg.Climates) {
			return
		}
		if !yield(DomainTypeNumber, reg.Numbers) {
			return
		}
		if !yield(DomainTypeDatetimeDate, reg.Dates) {
			return
		}
		if !yield(DomainTypeDatetimeTime, reg.Times) {
			return
		}
		if !yield(DomainTypeDatetimeDatetime, reg.Datetimes) {
			return
		}
		if !yield(DomainTypeText, reg.Texts) {
			return
		}
		if !yield(DomainTypeSelect, reg.Selects) {
			return
		}
		if !yield(DomainTypeLock, reg.Locks) {
			return
		}
		if !yield(DomainTypeValve, reg.Valves) {
			return
		}
		if !yield(DomainTypeMediaPlayer, reg.MediaPlayers) {
			return
		}
		if !yield(DomainTypeAlarmControlPanel, reg.AlarmControlPanels) {
			return
		}
		if !yield(DomainTypeEvent, reg.Events) {
			return
		}
		if !yield(DomainTypeUpdate, reg.Updates) {
			return
		}
	}
}

func IterateRegistry(reg *Registry) iter.Seq2[DomainType, Entity] {
	return func(yield func(DomainType, Entity) bool) {
		for t, f := range iteratorTypes(reg) {
			ents := f()
			for _, ent := range ents {
				if !yield(t, ent) {
					return
				}
			}
		}
	}
}
