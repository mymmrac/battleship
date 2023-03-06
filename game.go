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
	exit  bool

	connector BattleshipConnector
	events    chan Event

	currentScene *Scene
	scenes       map[SceneID]*Scene

	newGameBtn  *button
	joinGameBtn *button
	exitBtn     *button

	myBoard    *board
	myShipyard *shipyard

	readyBtn      *button
	notReadyBtn   *button
	clearBoardBtn *button

	opponentBoard *board

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

	buttonFace, err := loadFace(JetBrainsMonoFont, 24)
	if err != nil {
		return nil, err
	}

	newGameBtn := newButton(newPoint[float32](48, 48), 200, 40, "New Game", buttonFace)
	joinGameBtn := newButton(newPoint[float32](48, 48+40+32), 200, 40, "Join Game", buttonFace)
	exitBtn := newButton(newPoint[float32](48, 48+40*2+32*2), 200, 40, "Exit", buttonFace)

	boardFace, err := loadFace(JetBrainsMonoFont, float64(cellSize)*0.6)
	if err != nil {
		return nil, err
	}

	myBoard := newBoard(newPoint[float32](48, 48), boardFace)
	myShipyard := newShipyard(newPoint[float32](48, 440), myBoard, boardFace)
	opponentBoard := newBoard(newPoint[float32](48+400, 48), boardFace)

	readyBtn := newButton(newPoint[float32](48, 570), 120, 40, "Ready", buttonFace)
	notReadyBtn := newButton(newPoint[float32](48, 570), 160, 40, "Not Ready", buttonFace)
	clearBoardBtn := newButton(newPoint[float32](48+120+32, 570), 120, 40, "Clear", buttonFace)

	GlobalGameObjects.Acquire()
	defer GlobalGameObjects.Release()

	game := &Game{
		debug: false,

		connector: NewConnector(),
		events:    make(chan Event),

		newGameBtn:  RegisterObject(newGameBtn),
		joinGameBtn: RegisterObject(joinGameBtn),
		exitBtn:     RegisterObject(exitBtn),

		myBoard:    RegisterObject(myBoard),
		myShipyard: RegisterObject(myShipyard),

		readyBtn:      RegisterObject(readyBtn),
		notReadyBtn:   RegisterObject(notReadyBtn),
		clearBoardBtn: RegisterObject(clearBoardBtn),

		opponentBoard: RegisterObject(opponentBoard),

		objects: GlobalGameObjects.Objects(),
	}

	game.InitScenes()
	game.currentScene = game.scenes[sceneMenu]
	game.currentScene.OnEnter()

	return game, nil
}

func (g *Game) Update() error {
	if g.exit || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
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

	g.currentScene.OnUpdate()

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
