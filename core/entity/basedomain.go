package entity

import (
	"cmp"
	"fmt"
	"slices"
	"sync"

	"github.com/gosthome/gosthome/core/component"
)

func clone[T Entity](in []T) []Entity {
	ret := make([]Entity, 0, len(in))
	for _, e := range in {
		ret = append(ret, e)
	}
	return ret
}

func findByKey[T Entity](entityStore []T, hash uint32) (ret T, found bool) {
	i, found := slices.BinarySearchFunc(entityStore, hash, func(e T, t uint32) int {
		return cmp.Compare(e.HashID(), t)
	})
	if found {
		ret = entityStore[i]
	}
	return
}

func registerEntity[T Entity](entityStore []T, ent T) ([]T, error) {
	if ent.ID() == "" {
		return entityStore, fmt.Errorf("trying to register %T without an ID", ent)
	}
	i, found := slices.BinarySearchFunc(entityStore, ent, func(e T, t T) int {
		return cmp.Compare(e.HashID(), t.HashID())
	})
	if found {
		return entityStore, fmt.Errorf("hash id of %T is already registered %s!", ent, ent.ID())
	}
	return slices.Insert(entityStore, i, ent), nil
}

type BaseDomain[Domain any, EntityType EntityComponent, PD interface {
	DomainTyper
	*Domain
}] struct {
	mux sync.RWMutex

	entities []EntityType
}

func (bd *BaseDomain[Domain, EntityType, PD]) Setup() {
}

func (bd *BaseDomain[Domain, EntityType, PD]) InitializationPriority() component.InitializationPriority {
	return component.InitializationPriorityBus
}

func (bd *BaseDomain[Domain, EntityType, PD]) Close() error {
	bd.mux.Lock()
	defer bd.mux.Unlock()
	bd.entities = nil
	return nil
}

func (bd *BaseDomain[Domain, EntityType, PD]) Clone() []Entity {
	bd.mux.RLock()
	defer bd.mux.RUnlock()
	return clone(bd.entities)
}

func (bd *BaseDomain[Domain, EntityType, PD]) FindByKey(key uint32) (ret EntityType, found bool) {
	bd.mux.RLock()
	defer bd.mux.RUnlock()
	return findByKey(bd.entities, key)
}

func (bd *BaseDomain[Domain, EntityType, PD]) Register(ent EntityType) (err error) {
	bd.mux.Lock()
	defer bd.mux.Unlock()
	bd.entities, err = registerEntity(bd.entities, ent)
	return err
}

func (bd *BaseDomain[Domain, EntityType, PD]) DomainType() DomainType {
	return (PD)(nil).DomainType()
}
