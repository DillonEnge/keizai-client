package systems

import (
	"fmt"
	"math/rand"

	"github.com/DillonEnge/keizai-client/internal/components"
	"github.com/DillonEnge/keizai-client/internal/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

type WiggleSystem struct{}

func NewWiggleSystem() *WiggleSystem {
	return &WiggleSystem{}
}

func (w *WiggleSystem) Setup(e []*ecs.Entity) error {
	return nil
}

func (w *WiggleSystem) Draw(e []*ecs.Entity, screen *ebiten.Image) error {
	return nil
}

func (w *WiggleSystem) Update(e []*ecs.Entity) error {
	for _, v := range e {
		if v.Query(components.POSITION, components.IMAGE) {
			p, ok := v.Components[components.POSITION].(components.Position)
			if !ok {
				return fmt.Errorf("failed to cast position component")
			}

			x, y := rand.Intn(7), rand.Intn(7)

			p.X += x - 3
			p.Y += y - 3

			v.Components[components.POSITION] = p
		}
	}
	return nil
}
