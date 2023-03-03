package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	baseWindowWidth  = 1080
	baseWindowHeight = 720
)

type Game struct {
	debug bool
}

func NewGame() *Game {
	ebiten.SetWindowTitle("Battleship")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	deviceScale := ebiten.DeviceScaleFactor()
	windowWidth := baseWindowWidth * deviceScale
	windowHeight := baseWindowHeight * deviceScale
	ebiten.SetWindowSize(int(math.Ceil(windowWidth)), int(math.Ceil(windowHeight)))

	screenWidth, screenHeight := ebiten.ScreenSizeInFullscreen()
	ebiten.SetWindowPosition(
		int(math.Ceil((float64(screenWidth)-windowWidth)/2.0)),
		int(math.Ceil((float64(screenHeight)-windowHeight)/2.0)),
	)

	return &Game{
		debug: true,
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF10) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		g.debug = !g.debug
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.debug {
		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf("FPS: %0.2f\nTPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualTPS()),
			4, 4)

		cx, cy := ebiten.CursorPosition()
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("X: %d\nY: %d", cx, cy), cx+12, cy-40)

		size := screen.Bounds().Size()
		vector.StrokeLine(screen, 0, float32(cy), float32(size.X), float32(cy), 2, color.White)
		vector.StrokeLine(screen, float32(cx), 0, float32(cx), float32(size.Y), 2, color.White)
	}

	ebitenutil.DebugPrintAt(screen, "Hello World!", 100, 100)
}

func (g *Game) Layout(_, _ int) (int, int) {
	panic("unreachable")
}

func (g *Game) LayoutF(logicalWindowWidth, logicalWindowHeight float64) (float64, float64) {
	scale := ebiten.DeviceScaleFactor()
	return math.Ceil(logicalWindowWidth * scale), math.Ceil(logicalWindowHeight * scale)
}
