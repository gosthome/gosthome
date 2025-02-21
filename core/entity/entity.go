package entity

import (
	"hash/fnv"

	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/util"
	"github.com/oklog/ulid/v2"
)

//go:generate go-enum

// ENUM(ok,warning,error)
type EntityStatus uint8

type Entity interface {
	ID() string
	HashID() uint32
	Name() string
	Internal() bool
	DisabledByDefault() bool
	EntityCategory() EntityCategory
}

type EntityComponent interface {
	Entity
	component.Component
}

type WithIcon interface {
	Icon() string
}

type WithDeviceClass interface {
	DeviceClass() string
}

type WithUnitOfMeasurement interface {
	UnitOfMeasurement() string
}

type WithState[T any] interface {
	State() T
}

type BaseEntity struct {
	id       string
	idhash   uint32
	category EntityCategory
	name     string

	internal    bool
	idhashReady bool

	disabledByDefault bool
}

func MakeID(prefix string) string {
	return prefix + "_" + ulid.Make().String()
}

func NewBaseEntity(t DomainType, cfg *EntityConfig) *BaseEntity {
	b := &BaseEntity{
		id:                cfg.ID,
		name:              cfg.Name,
		internal:          cfg.Internal,
		category:          cfg.Category,
		disabledByDefault: cfg.DisabledByDefault,
	}
	if b.id == "" {
		if b.name != "" {
			b.id = util.CleanString(util.SnakeCase(b.name))
		} else {
			b.id = MakeID(t.String())
		}
	}
	b.idhash = HashID(b.id)
	return b
}

// ID implements Entity.
func (b *BaseEntity) ID() string {
	return b.id
}

func (b *BaseEntity) SetID(s string) {
	b.id = s
	b.idhash = HashID(b.id)
}

func HashID(id string) uint32 {
	h := fnv.New32()
	h.Write([]byte(id))
	return h.Sum32()
}

// HashID implements Entity.
func (b *BaseEntity) HashID() uint32 {
	if !b.idhashReady {
		b.idhash = HashID(b.id)
	}
	return b.idhash
}

// Name implements Entity.
func (b *BaseEntity) Name() string {
	return b.name
}

func (b *BaseEntity) SetName(name string) {
	b.name = name
}

// Internal implements Entity.
func (b *BaseEntity) Internal() bool {
	return b.internal
}

func (b *BaseEntity) SetInternal(internal bool) {
	b.internal = internal
}

// DisabledByDefault implements Entity.
func (b *BaseEntity) DisabledByDefault() bool {
	return b.disabledByDefault
}

func (b *BaseEntity) SetDisabledByDefault(disabledByDefault bool) {
	b.disabledByDefault = disabledByDefault
}

// EntityCategory implements Entity.
func (b *BaseEntity) EntityCategory() EntityCategory {
	return b.category
}

func (b *BaseEntity) SetEntityCategory(entityCategory EntityCategory) {
	b.category = entityCategory
}

var _ Entity = (*BaseEntity)(nil)

type IconMixin struct {
	icon string
}

func NewIconMixin(cfg *IconMixinConfig) IconMixin {
	return IconMixin{
		icon: cfg.Icon,
	}
}

func (b *IconMixin) SetIcon(icon string) {
	b.icon = icon
}

func (b *IconMixin) Icon() string {
	return b.icon
}

var _ WithIcon = (*IconMixin)(nil)

type DeviceClassMixin struct {
	deviceClass string
}

func NewDeviceClassMixin[Enum any, PE interface {
	DeviceClassValues
	*Enum
}](cfg *DeviceClassMixinConfig[Enum, PE]) DeviceClassMixin {
	return DeviceClassMixin{
		deviceClass: cfg.DeviceClass,
	}
}

func (b *DeviceClassMixin) SetDeviceClass(deviceClass string) {
	b.deviceClass = deviceClass
}

func (b *DeviceClassMixin) DeviceClass() string {
	return b.deviceClass
}

var _ WithDeviceClass = (*DeviceClassMixin)(nil)

type UnitOfMeasurementMixin struct {
	unitOfMeasurement string
}

func NewUnitOfMeasurementMixin(cfg *UnitOfMeasurementMixinConfig) UnitOfMeasurementMixin {
	return UnitOfMeasurementMixin{
		unitOfMeasurement: cfg.UnitOfMeasurement,
	}
}

func (b *UnitOfMeasurementMixin) SetUnitOfMeasurement(unitOfMeasurement string) {
	b.unitOfMeasurement = unitOfMeasurement
}

func (b *UnitOfMeasurementMixin) UnitOfMeasurement() string {
	return b.unitOfMeasurement
}

var _ WithUnitOfMeasurement = (*UnitOfMeasurementMixin)(nil)
