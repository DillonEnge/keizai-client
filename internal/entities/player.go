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
