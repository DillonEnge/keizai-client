package systems

import (
	"fmt"
	"math/rand"

	"github.com/DillonEnge/keizai-client/internal/components"
	"github.com/DillonEnge/keizai-client/internal/ecs"
)

type AiControllerSystem struct {
}

func NewAiControllerSystem() *AiControllerSystem {
	return &AiControllerSystem{}
}

func (a *AiControllerSystem) Update(e []*ecs.Entity) error {
	for _, v := range e {
		if v.Query(
			components.POSITION,
			components.IMAGE,
			components.SPEED,
			components.AI_CONTROLLER,
		) {
			p, ok := v.Components[components.POSITION].(components.Position)
			if !ok {
				return fmt.Errorf("failed to cast to position component")
			}

			s, ok := v.Components[components.SPEED].(components.Speed)
			if !ok {
				return fmt.Errorf("failed to cast to speed component")
			}

			a, ok := v.Components[components.AI_CONTROLLER].(components.AiController)
			if !ok {
				return fmt.Errorf("failed to cast to speed component")
			}

			switch rand.Intn(5) {
			case 0:
				a.UpdateCh <- func() {
					p.X += s.Speed
				}
			case 1:
				p.X -= s.Speed
			case 2:
				p.Y += s.Speed
			case 3:
				p.Y -= s.Speed
			default:
			}

			v.Components[components.POSITION] = p
		}
	}
	return nil
}
