package entities

import (
	"github.com/DillonEnge/keizai-client/internal/components"
	"github.com/DillonEnge/keizai-client/internal/ecs"
)

func NewNpcEntity(id *string, w, h, x, y, speed int, b components.BehaviorClass) *ecs.Entity {
	return ecs.NewEntity(
		id,
		components.NewImage(w, h),
		components.NewPosition(x, y),
		components.NewAiController(b),
		components.NewSpeed(speed),
	)
}
