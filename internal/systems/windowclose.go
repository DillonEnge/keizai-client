package systems

import (
	"fmt"
	"log/slog"

	"github.com/DillonEnge/keizai-client/internal/components"
	"github.com/DillonEnge/keizai-client/internal/ecs"
	keizai_grpc "github.com/DillonEnge/keizai-grpc"
	"github.com/hajimehoshi/ebiten/v2"
)

type WindowCloseSystem struct {
	Client keizai_grpc.KeizaiGrpcClient
}

func NewWindowCloseSystem(c keizai_grpc.KeizaiGrpcClient) *WindowCloseSystem {
	return &WindowCloseSystem{
		Client: c,
	}
}

func (w *WindowCloseSystem) Setup(e []*ecs.Entity) error {
	ebiten.SetWindowClosingHandled(true)

	return nil
}

func (w *WindowCloseSystem) Update(e []*ecs.Entity) error {
	if ebiten.IsWindowBeingClosed() {
		for _, v := range e {
			if v.Query(components.NETWORKED) {
				_, err := w.Client.RemoveEntity(v.Ctx, &keizai_grpc.RemoveEntityRequest{Id: v.Id})
				if err != nil {
					slog.Error("failed to remove entity", "err", err)
					return err
				}
			}
		}
		return fmt.Errorf("window closing")
	}

	return nil
}
