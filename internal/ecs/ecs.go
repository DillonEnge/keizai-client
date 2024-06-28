package ecs

import (
	"context"
	"fmt"

	"github.com/DillonEnge/keizai-client/internal/components"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type System interface {
	Setup([]*Entity) error
	Update([]*Entity) error
	Draw([]*Entity, *ebiten.Image) error
}

type Entity struct {
	Id         string
	Ctx        context.Context
	Cancel     context.CancelFunc
	Components map[components.Tag]interface{}
}

type Component interface{}

type TaggedComponent struct {
	Tag       components.Tag
	Component interface{}
}

func NewTaggedComponent(c Component) TaggedComponent {
	switch c.(type) {
	case components.Controllable:
		return TaggedComponent{
			Tag:       components.CONTROLLABLE,
			Component: c,
		}
	case components.Image:
		return TaggedComponent{
			Tag:       components.IMAGE,
			Component: c,
		}
	case components.Position:
		return TaggedComponent{
			Tag:       components.POSITION,
			Component: c,
		}
	case components.Speed:
		return TaggedComponent{
			Tag:       components.SPEED,
			Component: c,
		}
	case components.AiController:
		return TaggedComponent{
			Tag:       components.AI_CONTROLLER,
			Component: c,
		}
	case components.Network:
		return TaggedComponent{
			Tag:       components.NETWORKED,
			Component: c,
		}
	case components.Remote:
		return TaggedComponent{
			Tag:       components.REMOTE,
			Component: c,
		}
	default:
		panic(fmt.Errorf("failed to find tag for component %+v", c))
	}
}

func NewEntity(id *string, c ...Component) *Entity {
	if id == nil {
		u := uuid.NewString()
		id = &u
	}

	tc := []TaggedComponent{}
	for _, v := range c {
		tc = append(tc, NewTaggedComponent(v))
	}

	cm := make(map[components.Tag]interface{})
	for _, v := range tc {
		cm[v.Tag] = v.Component
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	return &Entity{
		Id:         *id,
		Ctx:        ctx,
		Cancel:     cancel,
		Components: cm,
	}
}

func (e *Entity) Query(tags ...components.Tag) bool {
	for _, v := range tags {
		if _, ok := e.Components[v]; !ok {
			return false
		}
	}

	return true
}

func (e *Entity) GetId() string {
	return e.Id
}

func (e *Entity) AddComponent(t TaggedComponent) error {
	_, ok := e.Components[t.Tag]
	if ok {
		return fmt.Errorf("component already present on entity")
	}

	e.Components[t.Tag] = t.Component

	return nil
}

func (e *Entity) RemoveComponent(tag components.Tag) error {
	_, ok := e.Components[tag]
	if !ok {
		return fmt.Errorf("component not present on entity")
	}

	delete(e.Components, tag)

	return nil
}

func GetComponent[T any](e Entity, tag components.Tag) (v T, ok bool) {
	v, ok = e.Components[tag].(T)
	return
}
