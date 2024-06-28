package statemanagement

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

type FSM struct {
	State State
}

type State interface {
	Update() error
	Draw(screen *ebiten.Image) error
	Teardown() error
	ProcessInput() error
}

func NewFSM(initialState State) *FSM {
	return &FSM{
		State: initialState,
	}
}

func (f *FSM) ChangeState(newState State) error {
	if err := f.State.Teardown(); err != nil {
		slog.Error("failed to teardown state", "err", err)
		return err
	}

	f.State = newState
	return nil
}
