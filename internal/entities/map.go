package entities

import (
	"github.com/DillonEnge/keizai-client/internal/components"
	"github.com/DillonEnge/keizai-client/internal/ecs"
)

func NewMapEntity(id *string, x, y, w, h int) *ecs.Entity {
	return ecs.NewEntity(
		id,
		components.NewImage(w, h),
		components.NewPosition(x, y),
		components.NewGrid(w, h),
	)
}
