package game

import (
	"context"
	"fmt"
	"image/color"
	"log/slog"

	"github.com/DillonEnge/keizai-client/internal/components"
	"github.com/DillonEnge/keizai-client/internal/ecs"
	"github.com/DillonEnge/keizai-client/internal/statemanagement"
	"github.com/DillonEnge/keizai-client/internal/states"
	keizai_grpc "github.com/DillonEnge/keizai-grpc"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	// "github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
)

type Game struct {
	Client       keizai_grpc.KeizaiGrpcClient
	Systems      []ecs.System
	Entities     []*ecs.Entity
	Ctx          context.Context
	UpdateCh     chan func()
	TextRenderer *etxt.Renderer
	statemanagement.FSM
}

func NewGame(
	ctx context.Context,
	client keizai_grpc.KeizaiGrpcClient,
	txtRenderer *etxt.Renderer,
	entities ...*ecs.Entity,
) *Game {
	es := []*ecs.Entity{}
	for _, e := range entities {
		es = append(es, e)
	}

	return &Game{
		Client:       client,
		Entities:     es,
		Ctx:          ctx,
		UpdateCh:     make(chan func(), 1),
		TextRenderer: txtRenderer,
		FSM:          *statemanagement.NewFSM(states.NewTitleState()),
	}
}

func (g *Game) AddSystems(s ...ecs.System) error {
	for _, v := range s {
		err := v.Setup(g.Entities)
		if err != nil {
			return err
		}

		g.Systems = append(g.Systems, v)
	}
	return nil
}

func (g *Game) ProcessInput() {
	if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		// g.FSM.ChangeState()
	}
	// k := make([]ebiten.Key, 0)
	// k = inpututil.AppendJustReleasedKeys(k)
	// fmt.Printf("%+v\n", k)
}

func (g *Game) Update() error {
	for _, v := range g.Systems {
		if err := v.Update(g.Entities); err != nil {
			return err
		}
	}

	if err := g.FSM.State.Update(); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%.0f TPS", ebiten.ActualTPS()))

	g.FSM.State.Draw(screen)

	for _, v := range g.Entities {
		if ok := v.Query(components.POSITION, components.IMAGE); ok {
			p, ok := v.Components[components.POSITION].(components.Position)
			if !ok {
				slog.Error("failed to cast component.", "id", v.Id)
			}

			d := ebiten.DrawImageOptions{}
			d.GeoM.Translate(
				float64(p.X),
				float64(p.Y),
			)
			i, ok := v.Components[components.IMAGE].(components.Image)
			if !ok {
				slog.Error("failed to cast component")
			}

			i.Image.Fill(color.White)

			screen.DrawImage(i.Image, &d)
			// slog.Info("drawing entity", "id", v.Id, "pos", v.Components[POSITION].(*PositionComponent))
		}
	}
	// for _, v := range g.Systems {
	// 	if err := v.Draw(g.Entities, screen); err != nil {
	// 		slog.Error("failed to draw system", "system", v)
	// 	}
	// }

	// g.TextRenderer.SetTarget(screen)
	// g.TextRenderer.SetColor(color.White)
	// g.TextRenderer.Draw("Hello World", 0, screen.Bounds().Dy())
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / 2, outsideHeight / 2
}
