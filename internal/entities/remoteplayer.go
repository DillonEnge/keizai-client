package entities

import (
	"github.com/DillonEnge/keizai-client/internal/components"
	"github.com/DillonEnge/keizai-client/internal/ecs"
)

func NewRemotePlayer(id string, x, y, w, h int) *ecs.Entity {
	return ecs.NewEntity(
		&id,
		components.NewPosition(x, y),
		components.NewImage(w, h),
		components.NewRemote(),
	)
}
