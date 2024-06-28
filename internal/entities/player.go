package entities

import (
	"github.com/DillonEnge/keizai-client/internal/components"
	"github.com/DillonEnge/keizai-client/internal/ecs"
)

func NewPlayerEntity(id *string, w, h, x, y, speed int) *ecs.Entity {
	return ecs.NewEntity(
		id,
		components.NewImage(w, h),
		components.NewPosition(x, y),
		components.NewControllable(),
		components.NewSpeed(speed),
		components.NewNetwork(),
	)
}

// import (
// 	"context"
// 	"image/color"
// 	"log/slog"

// 	"github.com/DillonEnge/keizai-client/internal/statemanagement"
// 	"github.com/DillonEnge/keizai-client/internal/stores"
// 	keizai_grpc "github.com/DillonEnge/keizai-grpc"
// 	"github.com/google/uuid"
// 	"github.com/hajimehoshi/ebiten/v2"
// )

// type Player struct {
// 	Store    stores.PlayerStore
// 	UpdateCh chan func()
// 	statemanagement.FSM
// }

// func NewPlayer(ctx context.Context, client keizai_grpc.KeizaiGrpcClient, id *string) *Player {
// 	if id == nil {
// 		s := uuid.NewString()
// 		id = &s
// 	}
// 	p := &Player{
// 		Store: *stores.NewPlayerStore(ctx, ebiten.NewImage(20, 20), client, *id, false)
// 		Image:         ebiten.NewImage(20, 20),
// 		Client:        client,
// 		Ctx:           ctx,
// 		Id:            *id,
// 		Authoritative: false,
// 		UpdateCh:      make(chan func(), 10),
// 	}
// }

// func (p *Player) SetAuthoritative(b bool) {
// 	p.Authoritative = b
// }

// func (p *Player) Update() error {
// 	if p.Authoritative {
// 		p.ProcessInput()
// 	}
// 	if len(p.UpdateCh) > 0 {
// 		u := <-p.UpdateCh
// 		u()
// 	}
// 	return nil
// }

// func (p *Player) ProcessInput() {
// 	newX, newY := p.Pos.X, p.Pos.Y
// 	if ebiten.IsKeyPressed(ebiten.KeyA) {
// 		newX--
// 		// p.Pos.X--
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyD) {
// 		newX++
// 		// p.Pos.X++
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyS) {
// 		newY++
// 		// p.Pos.Y++
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyW) {
// 		newY--
// 		// p.Pos.Y--
// 	}
// 	if newX != p.Pos.X || newY != p.Pos.Y {
// 		go p.Client.UpdatePosition(p.Ctx, &keizai_grpc.UpdatePositionRequest{
// 			Id:       p.Id,
// 			Position: &keizai_grpc.PositionComponent{X: newX, Y: newY},
// 		})
// 		if p.Authoritative {
// 			p.Pos.X, p.Pos.Y = newX, newY
// 		}
// 	}
// }

// func (p *Player) Draw(screen *ebiten.Image) {
// 	p.Image.Fill(color.RGBA{
// 		R: 255,
// 		G: 255,
// 		B: 255,
// 	})
// 	do := ebiten.DrawImageOptions{}
// 	do.GeoM.Translate(float64(p.Pos.X), float64(p.Pos.Y))
// 	screen.DrawImage(p.Image, &do)
// }

// func (p *Player) GetId() string {
// 	return p.Id
// }

// func (p *Player) AddEntity() {
// 	slog.Info("calling AddEntity", "id", p.Id)
// 	_, err := p.Client.AddEntity(p.Ctx, &keizai_grpc.AddEntityRequest{
// 		Id:       p.Id,
// 		Position: &keizai_grpc.PositionComponent{X: p.Pos.X, Y: p.Pos.Y},
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func (p *Player) StartPositionStream() {
// 	go func() {
// 		stream, err := p.Client.GetPosition(p.Ctx, &keizai_grpc.GetPositionRequest{Id: p.Id})
// 		if err != nil {
// 			panic(err)
// 		}
// 		for {
// 			select {
// 			case <-p.Ctx.Done():
// 				return
// 			default:
// 			}
// 			v, err := stream.Recv()
// 			if err != nil {
// 				slog.Error("failed to recv from stream", "err", err, "id", p.Id)
// 				break
// 			}
// 			p.UpdateCh <- func() {
// 				p.Pos = entities.Position{X: v.Position.X, Y: v.Position.Y}
// 			}
// 		}
// 	}()
// }
