package main

import (
	"context"
	"errors"
	"image/color"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"time"

	keizai_grpc "github.com/DillonEnge/keizai-grpc"
	"github.com/hajimehoshi/ebiten/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Game struct {
	boxes map[string]*Box
}

func NewGame() *Game {
	return &Game{
		boxes: make(map[string]*Box),
	}
}

func (g *Game) AddBox(box *Box) error {
	if _, ok := g.boxes[box.id]; ok {
		return errors.New("already present")
	}

	g.boxes[box.id] = box

	return nil
}

func (g *Game) DeleteBox(id string) error {
	if _, ok := g.boxes[id]; !ok {
		return errors.New("Box not found")
	}

	delete(g.boxes, id)

	return nil
}

type Setupable interface {
	Setup(context.Context) error
}

type Teardownable interface {
	Teardown(context.Context) error
}

type Updatable interface {
	Update(context.Context)
}

type Drawable interface {
	Draw(*ebiten.Image)
	GetPos() (x, y int32)
	GetDim() (w, h int32)
}

type Processable interface {
	Setupable
	Teardownable
	Drawable
	Updatable
}

type Box struct {
	grpcClient keizai_grpc.KeizaiGrpcClient
	pos        *keizai_grpc.PositionComponent
	dim        Dimensions
	id         string
	image      *ebiten.Image
	speed      int32
	controlled bool
	updateCh   chan *keizai_grpc.PositionComponent
	cancelFunc func()
}

func NewBox(grpcClient keizai_grpc.KeizaiGrpcClient, pos *keizai_grpc.PositionComponent, dim Dimensions, speed int32, controlled bool) *Box {
	return &Box{
		grpcClient: grpcClient,
		pos:        pos,
		dim:        dim,
		speed:      speed,
		image:      ebiten.NewImage(int(dim.Width), int(dim.Height)),
		controlled: controlled,
		updateCh:   make(chan *keizai_grpc.PositionComponent),
		cancelFunc: nil,
	}
}

type Dimensions struct {
	Width, Height int32
}

func (b *Box) GetPos() (x, y int32) {
	return b.pos.X, b.pos.Y
}

func (b *Box) GetDim() (w, h int32) {
	return b.dim.Width, b.dim.Height
}

func (b *Box) Setup(ctx context.Context) error {
	slog.Info("Setup called", "id", b.id)
	if !b.controlled {
		ctx, cancel := context.WithCancel(ctx)
		b.cancelFunc = cancel
		go schedule(ctx, func() {
			pos, err := b.grpcClient.GetPosition(ctx, &keizai_grpc.GetPositionRequest{
				Id: b.id,
			})
			if err != nil {
				slog.Error("Failed to get pos for remote box", "err", err)
				return
			}

			b.updateCh <- pos
		}, time.Millisecond*100)
		return nil
	}

	resp, err := b.grpcClient.CreateEntity(ctx, &keizai_grpc.CreateEntityRequest{
		Position: b.pos,
	})
	if err != nil {
		return err
	}

	b.id = resp.Id

	return nil
}

func (b *Box) Teardown(_ context.Context) error {
	slog.Info("Teardown called", "id", b.id)
	if !b.controlled {
		b.cancelFunc()
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := b.grpcClient.DeleteEntity(ctx, &keizai_grpc.DeleteEntityRequest{
		Id: b.id,
	})
	return err
}

func (b *Box) Draw(screen *ebiten.Image) {
	x, y := b.GetPos()
	d := ebiten.DrawImageOptions{}
	d.GeoM.Translate(float64(x), float64(y))
	b.image.Fill(color.Black)
	screen.DrawImage(b.image, &d)
}

func (b *Box) Update(ctx context.Context) {
	if !b.controlled {
		select {
		case newPos := <-b.updateCh:
			b.pos = newPos
		default:
		}
		return
	}

	pos := &keizai_grpc.PositionComponent{X: b.pos.X, Y: b.pos.Y}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		pos.X -= b.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		pos.X += b.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		pos.Y -= b.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		pos.Y += b.speed
	}

	if pos.X == b.pos.X && pos.Y == b.pos.Y {
		return
	}

	slog.Info("processed new position", "pos", pos)

	b.pos = pos

	b.grpcClient.UpdatePosition(ctx, &keizai_grpc.UpdatePositionRequest{
		Id:       b.id,
		Position: b.pos,
	})
}

func (g *Game) Update() error {
	for _, v := range g.boxes {
		v.Update(context.TODO())
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	for _, v := range g.boxes {
		v.Draw(screen)
	}
	// ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	client, err := grpc.NewClient("keizai-server.engehost.net:8765", opts...)
	if err != nil {
		return err
	}

	keizaiGrpcClient := keizai_grpc.NewKeizaiGrpcClient(client)

	b := NewBox(
		keizaiGrpcClient,
		&keizai_grpc.PositionComponent{X: 100, Y: 100},
		Dimensions{Width: 20, Height: 20},
		2,
		true,
	)

	b.Setup(ctx)
	defer func() {
		if err := b.Teardown(ctx); err != nil {
			slog.Error("Failed to teardown box", "id", b.id, "err", err)
		}
	}()

	g := NewGame()

	g.AddBox(b)

	// for _, v := range w.boxes {
	// 	if err := v.Setup(ctx); err != nil {
	// 		slog.Error("Failed to setup box", "err", err)
	// 		continue
	// 	}

	// 	defer func() {
	// 		if err := v.Teardown(ctx); err != nil {
	// 			slog.Error("Failed to teardown box", "id", v.id, "err", err)
	// 		}
	// 	}()
	// }

	go schedule(ctx, func() {
		resp, err := keizaiGrpcClient.GetEntityIds(ctx, nil)
		if err != nil {
			slog.Error("Failed to get entity ids", "err", err)
			return
		}

		for k := range g.boxes {
			if !slices.Contains(resp.Ids, k) {
				g.boxes[k].Teardown(ctx)
				g.DeleteBox(k)
			}
		}

		for _, v := range resp.Ids {
			nb := NewBox(keizaiGrpcClient, &keizai_grpc.PositionComponent{}, Dimensions{Width: 40, Height: 40}, 2, false)
			nb.id = v

			if err := g.AddBox(nb); err != nil {
				continue
			}

			if err := nb.Setup(ctx); err != nil {
				slog.Error("Failed to setup box", "err", err)
			}
		}
	}, time.Millisecond*200)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	if err := ebiten.RunGame(g); err != nil {
		return err
	}

	return nil
}

func schedule(ctx context.Context, fn func(), duration time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(duration):
			fn()
		}
	}
}