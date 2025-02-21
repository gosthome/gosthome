package registry

import (
	"fmt"
	"maps"
	"sync"

	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/entity"
)

type componentDeclMap = map[string]component.Declaration
type entityComponentMap = map[entity.DomainType]componentDeclMap

type Registry struct {
	reg   componentDeclMap
	ecReg entityComponentMap
}

func NewRegistry() *Registry {
	ecReg := make(entityComponentMap)
	for dt := range entity.DomainTypeEnd + 1 {
		ecReg[dt] = make(componentDeclMap)
	}
	return &Registry{
		reg:   make(componentDeclMap),
		ecReg: ecReg,
	}
}

func (cr *Registry) Get(name string) (cd component.Declaration, ok bool) {
	cd, ok = cr.reg[name]
	return
}

func (cr *Registry) GetEntityComponent(domain entity.DomainType, platform string) (cd component.Declaration, ok bool) {
	dr, ok := cr.ecReg[domain]
	if !ok {
		return nil, false
	}
	cd, ok = dr[platform]
	return
}

func (cr *Registry) Register(name string, cd component.Declaration) error {
	_, ok := cr.reg[name]
	if ok {
		return fmt.Errorf("component %s already registered", name)
	}
	cr.reg[name] = cd
	if err := cr.tryRegisterPlatforms(name, cd); err != nil {
		return err
	}
	return nil
}

func (cr *Registry) RegisterEntityComponent(domain entity.DomainType, platform string, cd component.Declaration) error {
	dr, ok := cr.ecReg[domain]
	if !ok {
		panic("Unknown domain or uninitialized Registry")
	}
	_, ok = dr[platform]
	if ok {
		return fmt.Errorf("%s component for platform %s already registered", domain.String(), platform)
	}
	dr[platform] = cd
	cr.ecReg[domain] = dr
	return nil
}

var (
	defaultRegistry    = NewRegistry()
	defaultRegistryMux = sync.Mutex{}
)

func DefaultRegistry() *Registry {
	ret := NewRegistry()
	defaultRegistryMux.Lock()
	defer defaultRegistryMux.Unlock()
	ret.reg = maps.Clone(defaultRegistry.reg)
	ret.ecReg = maps.Clone(defaultRegistry.ecReg)
	return ret
}

func RegisterDefaultComponent(name string, cd component.Declaration) byte {
	defaultRegistryMux.Lock()
	defer defaultRegistryMux.Unlock()
	err := defaultRegistry.Register(name, cd)
	if err != nil {
		panic(err)
	}
	return 0
}

func RegisterDefaultEntityComponent(domain entity.DomainType, platform string, cd component.Declaration) byte {
	defaultRegistryMux.Lock()
	defer defaultRegistryMux.Unlock()
	err := defaultRegistry.RegisterEntityComponent(domain, platform, cd)
	if err != nil {
		panic(err)
	}
	return 0
}
