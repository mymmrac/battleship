package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	baseWindowWidth  = 1080
	baseWindowHeight = 720

	dpi = 72
)

type Game struct {
	debug         bool
	myBoard       *board
	opponentBoard *board

	hoverOnCell bool
	hoverCell   point[int]
}

func NewGame() (*Game, error) {
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

	boardFont, err := LoadFont("JetBrainsMono-Regular.ttf")
	if err != nil {
		return nil, fmt.Errorf("load font: %w", err)
	}

	boardFace, err := opentype.NewFace(boardFont, &opentype.FaceOptions{
		Size:    float64(cellSize) * 0.6,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("create font face: %w", err)
	}

	return &Game{
		debug:         true,
		myBoard:       newBoard(newPoint[float32](48, 48), boardFace),
		opponentBoard: newBoard(newPoint[float32](48+400, 48), boardFace),
	}, nil
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

	cx, cy := ebiten.CursorPosition()
	cp := newPoint(float32(cx), float32(cy))

	g.hoverCell.x, g.hoverCell.y, g.hoverOnCell = g.opponentBoard.cellOn(cp)
	if g.hoverOnCell && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		_ = g.opponentBoard.shoot(g.hoverCell.x, g.hoverCell.y)
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

	g.myBoard.draw(screen)
	g.opponentBoard.draw(screen)

	if g.hoverOnCell {
		pos := g.opponentBoard.cellPos(g.hoverCell.x+1, g.hoverCell.y+1)
		vector.StrokeRect(screen, pos.x, pos.y, cellSize, cellSize, 4, color.RGBA{
			R: 236, // 149,
			G: 168, // 189,
			B: 105, // 255,
			A: 255,
		})
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	panic("unreachable")
}

func (g *Game) LayoutF(logicalWindowWidth, logicalWindowHeight float64) (float64, float64) {
	scale := ebiten.DeviceScaleFactor()
	return math.Ceil(logicalWindowWidth * scale), math.Ceil(logicalWindowHeight * scale)
}
