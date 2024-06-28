package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ChattingState struct {
	Message string
}

func NewChattingState() *ChattingState {
	return &ChattingState{}
}

func (c *ChattingState) Update() error {
	if err := c.ProcessInput(); err != nil {
		panic(err)
	}
	return nil
}

func (c *ChattingState) Teardown() error {
	return nil
}

func (c *ChattingState) ProcessInput() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// g.FSM.ChangeState(&PlayingState{})
		return nil
	}

	k := make([]ebiten.Key, 0)
	k = inpututil.AppendJustPressedKeys(k)
	for _, k := range k {
		c.Message += k.String()
	}
	return nil
}
