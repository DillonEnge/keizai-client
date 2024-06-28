package systems

import (
	"fmt"
	"image/color"

	"github.com/DillonEnge/keizai-client/internal/components"
	"github.com/DillonEnge/keizai-client/internal/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

type RenderSystem struct {
	Entities *[]*ecs.Entity
}

func NewRenderSystem() *RenderSystem {
	return &RenderSystem{}
}

func (r *RenderSystem) Setup(e []*ecs.Entity) error {
	return nil
}

func (r *RenderSystem) Update(e []*ecs.Entity) error {
	return nil
}

func (r *RenderSystem) Draw(e []*ecs.Entity, screen *ebiten.Image) error {
	for _, v := range e {
		if ok := v.Query(components.POSITION, components.IMAGE); ok {
			p, ok := v.Components[components.POSITION].(components.Position)
			if !ok {
				return fmt.Errorf("failed to cast component. System: %s, Component: %s", "RenderSystem", "position")
			}

			d := ebiten.DrawImageOptions{}
			d.GeoM.Translate(
				float64(p.X),
				float64(p.Y),
			)
			i, ok := v.Components[components.IMAGE].(components.Image)
			if !ok {
				return fmt.Errorf("failed to cast component")
			}

			i.Image.Fill(color.White)

			screen.DrawImage(i.Image, &d)
			// slog.Info("drawing entity", "id", v.Id, "pos", v.Components[POSITION].(*PositionComponent))
		}
	}
	return nil
}
