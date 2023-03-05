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

	state gameState

	myBoard       *board
	myShipyard    *shipyard
	opponentBoard *board
	readyBtn      *button
	notReadyBtn   *button
	clearBoardBtn *button

	objects []GameObject
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

	readyBtn := newButton(newPoint[float32](48, 570), 120, 40, "Ready", buttonFace)
	readyBtn.Disable()

	notReadyBtn := newButton(newPoint[float32](48, 570), 160, 40, "Not Ready", buttonFace)
	notReadyBtn.Disable()
	notReadyBtn.Hide()

	clearBoardBtn := newButton(newPoint[float32](48+120+32, 570), 120, 40, "Clear", buttonFace)

	GlobalGameObjects.Acquire()
	defer GlobalGameObjects.Release()

	return &Game{
		debug: false,

		state: statePlaceShips,

		myBoard:       RegisterObject(myBoard),
		myShipyard:    RegisterObject(myShipyard),
		opponentBoard: RegisterObject(opponentBoard),
		readyBtn:      RegisterObject(readyBtn),
		notReadyBtn:   RegisterObject(notReadyBtn),
		clearBoardBtn: RegisterObject(clearBoardBtn),

		objects: GlobalGameObjects.Objects(),
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

	cursorPointer := false

	for _, updatable := range g.objects {
		if updatable.Active() {
			updatable.Update(cp)

			if updatable.CursorPointer() {
				cursorPointer = true
			}
		}
	}

	if cursorPointer {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	switch g.state {
	case statePlaceShips:
		if g.myBoard.hover {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				g.myBoard.placeShip(g.myBoard.hoverX, g.myBoard.hoverY)
			}

			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
				g.myBoard.removeShip(g.myBoard.hoverX, g.myBoard.hoverY)
			}
		}

		if g.clearBoardBtn.clicked {
			g.myBoard.cells = [cellsCount][cellsCount]cellKind{}
		}

		if g.myShipyard.ready() {
			g.readyBtn.Enable()

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

		if g.clearBoardBtn.clicked {
			g.myBoard.cells = [cellsCount][cellsCount]cellKind{}
		}

		if !g.myShipyard.ready() {
			g.readyBtn.Disable()

			g.clearBoardBtn.Enable()
			g.clearBoardBtn.Show()

			g.state = statePlaceShips
		}

		if g.readyBtn.clicked {
			g.readyBtn.Disable()
			g.readyBtn.Hide()

			g.notReadyBtn.Enable()
			g.notReadyBtn.Show()

			g.clearBoardBtn.Disable()
			g.clearBoardBtn.Hide()

			g.myBoard.Disable()
			g.opponentBoard.Enable()

			g.state = statePlayerReady
		}

	case statePlayerReady:
		if g.opponentBoard.hover && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			_ = g.opponentBoard.shoot(g.opponentBoard.hoverX, g.opponentBoard.hoverY)
		}

		if g.notReadyBtn.clicked {
			g.readyBtn.Enable()
			g.readyBtn.Show()

			g.notReadyBtn.Disable()
			g.notReadyBtn.Hide()

			g.clearBoardBtn.Enable()
			g.clearBoardBtn.Show()

			g.myBoard.Enable()
			g.opponentBoard.Disable()

			g.state = stateShipsPlaced
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

	for _, drawable := range g.objects {
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
