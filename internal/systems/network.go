package systems

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/DillonEnge/keizai-client/internal/components"
	"github.com/DillonEnge/keizai-client/internal/ecs"
	"github.com/DillonEnge/keizai-client/internal/entities"
	"github.com/DillonEnge/keizai-client/internal/game"
	keizai_grpc "github.com/DillonEnge/keizai-grpc"
)

type Network struct {
	Client   keizai_grpc.KeizaiGrpcClient
	Game     *game.Game
	Ctx      context.Context
	UpdateCh chan func() error
}

func NewNetwork(ctx context.Context, c keizai_grpc.KeizaiGrpcClient, g *game.Game) *Network {
	return &Network{
		Client:   c,
		Game:     g,
		Ctx:      ctx,
		UpdateCh: make(chan func() error, 10),
	}
}

func (n *Network) Setup(e []*ecs.Entity) error {
	for _, v := range e {
		if v.Query(components.POSITION, components.NETWORKED) {
			p, ok := v.Components[components.POSITION].(components.Position)
			if !ok {
				return fmt.Errorf("failed to cast to position component")
			}
			slog.Info("attempting to add entity")
			go n.Client.AddEntity(v.Ctx, &keizai_grpc.AddEntityRequest{
				Id:       v.Id,
				Position: &keizai_grpc.PositionComponent{X: int32(p.X), Y: int32(p.Y)},
			})
		}
		if v.Query(components.POSITION, components.REMOTE) {
			go n.StartGetPositionStream(v)
		}
	}
	go n.StartGetEntityIdsStream()
	return nil
}

func (n *Network) Update(e []*ecs.Entity) error {
	if len(n.UpdateCh) > 0 {
		update := <-n.UpdateCh
		if err := update(); err != nil {
			return err
		}
	}

	for _, v := range e {
		if v.Query(components.POSITION, components.NETWORKED) {
			p, ok := v.Components[components.POSITION].(components.Position)
			if !ok {
				slog.Error("failed to cast to position component")
				return fmt.Errorf("failed to cast to position component")
			}
			go n.Client.UpdatePosition(context.Background(), &keizai_grpc.UpdatePositionRequest{
				Id:       v.Id,
				Position: &keizai_grpc.PositionComponent{X: int32(p.X), Y: int32(p.Y)},
			})
		}
	}
	return nil
}

func (n *Network) StartGetEntityIdsStream() {
	stream, err := n.Client.GetEntityIds(context.Background(), nil)
	if err != nil {
		slog.Error("failed to start GetEntityIds stream", "err", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			slog.Error("failed to recv from stream", "err", err)
			break
		}
		for _, v := range res.Ids {
			if !slices.ContainsFunc(n.Game.Entities, func(entity *ecs.Entity) bool {
				return entity.Id == v
			}) {
				n.UpdateCh <- func() error {
					r := entities.NewRemotePlayer(v, 10, 10, 10, 10)

					go n.StartGetPositionStream(r)
					n.Game.Entities = append(n.Game.Entities, r)
					return nil
				}
			}
		}
		for i, v := range n.Game.Entities {
			if v.Query(components.REMOTE) {
				if !slices.ContainsFunc(res.Ids, func(id string) bool {
					return id == v.Id
				}) {
					n.UpdateCh <- func() error {
						v.Cancel()
						n.Game.Entities = append(n.Game.Entities[:i], n.Game.Entities[i+1:]...)
						return nil
					}
				}
			}
		}
	}
}

func (n *Network) StartGetPositionStream(e *ecs.Entity) {
	stream, err := n.Client.GetPosition(e.Ctx, &keizai_grpc.GetPositionRequest{Id: e.Id})
	if err != nil {
		slog.Error("failed to start GetPosition stream", "err", err)
		return
	}
	for {
		select {
		case <-e.Ctx.Done():
			return
		default:
		}
		v, err := stream.Recv()
		if err != nil {
			slog.Error("failed to recv from stream", "err", err, "id", e.Id)
			break
		}
		n.UpdateCh <- func() error {
			p, ok := e.Components[components.POSITION].(components.Position)
			if !ok {
				return fmt.Errorf("failed to cast to position component")
			}
			p.X, p.Y = int(v.Position.X), int(v.Position.Y)
			e.Components[components.POSITION] = p
			return nil
		}
	}
}
