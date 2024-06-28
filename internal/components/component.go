package components

import "github.com/hajimehoshi/ebiten/v2"

type Tag int

const (
	POSITION Tag = iota
	IMAGE
	CONTROLLABLE
	SPEED
	NETWORKED
	AI_CONTROLLER
	REMOTE
)

type Position struct {
	X int
	Y int
}

func NewPosition(x, y int) Position {
	return Position{
		X: x,
		Y: y,
	}
}

type Image struct {
	Image *ebiten.Image
}

func NewImage(w, h int) Image {
	return Image{
		Image: ebiten.NewImage(w, h),
	}
}

type Controllable struct{}

func NewControllable() Controllable {
	return Controllable{}
}

type Speed struct {
	Speed int
}

func NewSpeed(speed int) Speed {
	return Speed{
		Speed: speed,
	}
}

type Network struct {
	UpdateCh chan func()
}

func NewNetwork() Network {
	return Network{
		UpdateCh: make(chan func(), 1),
	}
}

type BehaviorClass int

const (
	NPC BehaviorClass = iota
)

type AiController struct {
	BehaviorClass BehaviorClass
	UpdateCh      chan func()
}

func NewAiController(b BehaviorClass) AiController {
	return AiController{
		BehaviorClass: b,
		UpdateCh:      make(chan func()),
	}
}

type Remote struct{}

func NewRemote() Remote {
	return Remote{}
}
