package systems

import (
	"fmt"

	"github.com/DillonEnge/keizai-client/internal/components"
	"github.com/DillonEnge/keizai-client/internal/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

type ControllerSystem struct{}

func NewControllerSystem() *ControllerSystem {
	return &ControllerSystem{}
}

func (c *ControllerSystem) Setup(e []*ecs.Entity) error {
	return nil
}

func (c *ControllerSystem) Update(e []*ecs.Entity) error {
	for _, v := range e {
		if v.Query(components.POSITION, components.IMAGE, components.CONTROLLABLE, components.SPEED) {
			c.ProcessInput(v)
		}
	}
	return nil
}

func (c *ControllerSystem) ProcessInput(e *ecs.Entity) error {
	p, ok := e.Components[components.POSITION].(components.Position)
	if !ok {
		return fmt.Errorf("failed to cast position component")
	}
	s, ok := e.Components[components.SPEED].(components.Speed)
	if !ok {
		return fmt.Errorf("failed to cast speed component")
	}
	newX, newY := p.X, p.Y
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		newX -= s.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		newX += s.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		newY += s.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		newY -= s.Speed
	}
	if newX != p.X || newY != p.Y {
		// go p.Client.UpdatePosition(p.Ctx, &keizai_grpc.UpdatePositionRequest{
		// 	Id:       p.Id,
		// 	Position: &keizai_grpc.PositionComponent{X: newX, Y: newY},
		// })
		// if p.Authoritative {
		p.X, p.Y = newX, newY
		// }
		e.Components[components.POSITION] = p
	}

	return nil
}
