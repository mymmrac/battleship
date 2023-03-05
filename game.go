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

type gameState int

const (
	statePlaceShips gameState = iota
	stateShipsPlaced
	statePlayerReady
)

type Game struct {
	debug bool

	updatable []Updatable
	drawable  []Drawable

	state gameState

	myBoard    *board
	myShipyard *shipyard

	opponentBoard *board

	playerReadyBtn *button
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

	boardFace, err := loadFace(JetBrainsMonoFont, float64(cellSize)*0.6)
	if err != nil {
		return nil, err
	}

	buttonFace, err := loadFace(JetBrainsMonoFont, 24)
	if err != nil {
		return nil, err
	}

	myBoard := newBoard(newPoint[float32](48, 48), boardFace)
	myShipyard := newShipyard(newPoint[float32](48, 440), myBoard, boardFace)

	opponentBoard := newBoard(newPoint[float32](48+400, 48), boardFace)
	opponentBoard.Disable()

	playerReadyBtn := newButton(newPoint[float32](48, 570), 120, 40, "Ready", buttonFace)
	playerReadyBtn.Disable()

	return &Game{
		debug: false,

		updatable: []Updatable{
			myBoard,
			opponentBoard,
			myShipyard,
			playerReadyBtn,
		},

		drawable: []Drawable{
			myBoard,
			opponentBoard,
			myShipyard,
			playerReadyBtn,
		},

		myBoard:    myBoard,
		myShipyard: myShipyard,

		opponentBoard: opponentBoard,

		playerReadyBtn: playerReadyBtn,
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

	for _, updatable := range g.updatable {
		if updatable.Active() {
			updatable.Update(cp)
		}
	}

	switch g.state {
	case statePlaceShips:
		OnHover(g.myBoard.hover)

		if g.myBoard.hover {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				g.myBoard.placeShip(g.myBoard.hoverX, g.myBoard.hoverY)
			}

			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
				g.myBoard.removeShip(g.myBoard.hoverX, g.myBoard.hoverY)
			}
		}

		if g.myShipyard.ready() {
			g.playerReadyBtn.Enable()

			g.state = stateShipsPlaced
		}
	case stateShipsPlaced:
		if g.myBoard.hover {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				g.myBoard.placeShip(g.myBoard.hoverX, g.myBoard.hoverY)
			}

			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
				g.myBoard.removeShip(g.myBoard.hoverX, g.myBoard.hoverY)
			}
		}

		if !g.myShipyard.ready() {
			g.playerReadyBtn.Disable()

			g.state = statePlaceShips
		}

		OnHover(g.myBoard.hover || g.playerReadyBtn.hover)

		if g.playerReadyBtn.clicked {
			g.playerReadyBtn.Disable()
			g.playerReadyBtn.Hide()

			g.myBoard.Disable()
			g.opponentBoard.Enable()

			g.state = statePlayerReady
		}
	case statePlayerReady:
		OnHover(g.opponentBoard.hover)
		if g.opponentBoard.hover && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			_ = g.opponentBoard.shoot(g.opponentBoard.hoverX, g.opponentBoard.hoverY)
		}
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

	for _, drawable := range g.drawable {
		if drawable.Visible() {
			drawable.Draw(screen)
		}
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	panic("unreachable")
}

func (g *Game) LayoutF(logicalWindowWidth, logicalWindowHeight float64) (float64, float64) {
	scale := ebiten.DeviceScaleFactor()
	return math.Ceil(logicalWindowWidth * scale), math.Ceil(logicalWindowHeight * scale)
}

type Updatable interface {
	Active() bool
	Update(cp point[float32])
}

type Drawable interface {
	Visible() bool
	Draw(screen *ebiten.Image)
}

func OnHover(hover bool) {
	if hover {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
}
