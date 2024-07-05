package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"

	// "log/slog"
	"time"

	"github.com/DillonEnge/keizai-client/internal/cert"
	"github.com/DillonEnge/keizai-client/internal/ecs"
	"github.com/DillonEnge/keizai-client/internal/entities"
	"github.com/DillonEnge/keizai-client/internal/fonts"
	"github.com/DillonEnge/keizai-client/internal/game"
	"github.com/DillonEnge/keizai-client/internal/systems"
	"github.com/DillonEnge/keizai-grpc"
	"github.com/hajimehoshi/ebiten/v2"

	// "github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
	// "github.com/tinne26/fonts/liberation/lbrtserif"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ChatBox struct {
	Text    string
	TimeOut time.Duration
}

func NewChatBox(text string, timeOut time.Duration) *ChatBox {
	return &ChatBox{
		Text:    text,
		TimeOut: timeOut,
	}
}

func (c *ChatBox) Open() {

}

type ImageConfig struct {
	Image    *ebiten.Image
	Position entities.Position
	Rotation float64
}

var (
	debug       = flag.Bool("d", false, "enable debug")
	entityCount = flag.Int("e", 1, "set number of entitites to render")
)

func main() {
	flag.Parse()
	r := bytes.NewReader(cert.C)
	creds, err := newClientTLS(r, "keizai-server.engehost.net")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(creds))

	conn, err := grpc.NewClient("keizai-server.engehost.net:8765", opts...)
	if err != nil {
		log.Fatalf("failed to create new client: %v", err)
	}
	defer conn.Close()

	client := keizai_grpc.NewKeizaiGrpcClient(conn)

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("KeiZai")

	// create a new text renderer and configure it
	txtRenderer, err := newTxtRenderer()
	if err != nil {
		log.Fatal("Error while loading fonts: " + err.Error())
	}

	es := []*ecs.Entity{
		entities.NewPlayerEntity(nil, 15, 15, 20, 20, 2),
	}

	ctx := context.Background()

	g := game.NewGame(ctx, client, txtRenderer, es...)

	ebiten.SetWindowClosingHandled(true)

	s := []ecs.System{
		// systems.NewRenderSystem(),
		systems.NewControllerSystem(),
		systems.NewNetwork(ctx, client, g),
		systems.NewWindowCloseSystem(client),
	}

	err = g.AddSystems(s...)
	if err != nil {
		slog.Error("failed to add systems", "err", err)
	}

	if err := ebiten.RunGame(g); err != nil {
		slog.Error("game terminated", "err", err)
		log.Fatal(err)
	}
}

func newClientTLS(r io.Reader, serverNameOverride string) (credentials.TransportCredentials, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	return credentials.NewTLS(&tls.Config{ServerName: serverNameOverride, RootCAs: cp}), nil
}

func newTxtRenderer() (*etxt.Renderer, error) {
	robotoFont := fonts.T

	fontLib := etxt.NewFontLibrary()
	_, err := fontLib.ParseFontBytes(robotoFont)
	// _, _, err = fontLib.ParseDirFonts("assets/fonts")
	if err != nil {
		return nil, err
	}

	// create a new text renderer and configure it
	txtRenderer := etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	txtRenderer.SetCacheHandler(glyphsCache.NewHandler())
	txtRenderer.SetFont(fontLib.GetFont("Roboto"))
	txtRenderer.SetAlign(etxt.Bottom, etxt.Left)
	txtRenderer.SetSizePx(12)

	return txtRenderer, nil
}
