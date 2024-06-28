package states

import (
	"github.com/DillonEnge/keizai-client/internal/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

type TitleState struct {
	Entities []ecs.Entity
}

func NewTitleState(entities ...ecs.Entity) *TitleState {
	return &TitleState{
		Entities: entities,
	}
}

func (t *TitleState) Update() error {

	return nil
}

func (t *TitleState) ProcessInput() error {

	return nil
}

func (t *TitleState) Teardown() error {
	return nil
}

func (t *TitleState) Draw(screen *ebiten.Image) error {

	return nil
}
