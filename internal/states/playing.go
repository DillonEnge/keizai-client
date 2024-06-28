package states

import "github.com/DillonEnge/keizai-client/internal/ecs"

type PlayingState struct {
	Entities []ecs.Entity
}

func NewPlayingState(entities ...ecs.Entity) *PlayingState {
	return &PlayingState{}
}

func (p *PlayingState) Update() error {
	return nil
}

func (p *PlayingState) TearDown() error {
	return nil
}

func (p *PlayingState) ProcessInput() error {
	return nil
}
