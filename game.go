package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"google.golang.org/grpc"

	"github.com/mymmrac/battleship/data"
	"github.com/mymmrac/battleship/events"
	"github.com/mymmrac/battleship/ui"
)

const (
	baseWindowWidth  = 1080
	baseWindowHeight = 720
)

type Game struct {
	debug bool
	exit  bool

	serverAddr string
	serverPort string

	grpcConn     *grpc.ClientConn
	eventManager *EventManagerClient

	events chan events.GameEvent

	currentScene *Scene
	scenes       map[SceneID]*Scene

	newGameBtn  *ui.Button
	joinGameBtn *ui.Button
	exitBtn     *ui.Button

	newGameLoadingLabel *ui.Label

	myBoard    *Board
	myShipyard *Shipyard

	readyBtn      *ui.Button
	notReadyBtn   *ui.Button
	clearBoardBtn *ui.Button

	opponentReady      bool
	opponentReadyLabel *ui.Label

	myTurn          bool
	lastShootPos    data.Point[int]
	playerTurnLabel *ui.Label
	opponentBoard   *Board

	theEndLabel *ui.Label

	objects []GameObject
}

func NewGame(serverAddr, serverPort string) (*Game, error) {
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

	i16, err := LoadImage(Icon16x16)
	if err != nil {
		return nil, err
	}
	i32, err := LoadImage(Icon32x32)
	if err != nil {
		return nil, err
	}
	i48, err := LoadImage(Icon48x48)
	if err != nil {
		return nil, err
	}
	ebiten.SetWindowIcon([]image.Image{i16, i32, i48})

	buttonFace, err := loadFace(JetBrainsMonoFont, 24)
	if err != nil {
		return nil, err
	}

	labelFace, err := loadFace(JetBrainsMonoFont, 32)
	if err != nil {
		return nil, err
	}

	newGameBtn := ui.NewButton(data.NewPoint[float32](48, 48), 200, 40, "New Game", buttonFace)
	joinGameBtn := ui.NewButton(data.NewPoint[float32](48, 48+40+32), 200, 40, "Join Game", buttonFace)
	exitBtn := ui.NewButton(data.NewPoint[float32](48, 48+40*2+32*2), 200, 40, "Exit", buttonFace)

	newGameLoadingLabel := ui.NewLabel(data.NewPoint[float32](48, 48), "", labelFace)

	boardFace, err := loadFace(JetBrainsMonoFont, float64(cellSize)*0.6)
	if err != nil {
		return nil, err
	}

	myBoard := NewBoard(data.NewPoint[float32](48, 48), boardFace)
	myShipyard := NewShipyard(data.NewPoint[float32](48, 440), myBoard, boardFace)
	opponentBoard := NewBoard(data.NewPoint[float32](48+400, 48), boardFace)

	readyBtn := ui.NewButton(data.NewPoint[float32](48, 570), 120, 40, "Ready", buttonFace)
	notReadyBtn := ui.NewButton(data.NewPoint[float32](48, 570), 160, 40, "Not Ready", buttonFace)
	clearBoardBtn := ui.NewButton(data.NewPoint[float32](48+120+32, 570), 120, 40, "Clear", buttonFace)
	opponentReadyLabel := ui.NewLabel(data.NewPoint[float32](48, 640), "Opponent: not ready", labelFace)

	playerTurnLabel := ui.NewLabel(data.NewPoint[float32](48+400+48/2, 440), "", labelFace)
	playerTurnLabel.SetAlignment(ui.LabelAlignmentTopCenter)

	theEndLabel := ui.NewLabel(data.NewPoint[float32](48+400+48/2, 440), "", labelFace)
	theEndLabel.SetAlignment(ui.LabelAlignmentTopCenter)

	GlobalGameObjects.Acquire()
	defer GlobalGameObjects.Release()

	game := &Game{
		debug: false,

		serverAddr: serverAddr,
		serverPort: serverPort,

		events: make(chan events.GameEvent),

		newGameBtn:  RegisterObject(newGameBtn),
		joinGameBtn: RegisterObject(joinGameBtn),
		exitBtn:     RegisterObject(exitBtn),

		newGameLoadingLabel: RegisterObject(newGameLoadingLabel),

		myBoard:    RegisterObject(myBoard),
		myShipyard: RegisterObject(myShipyard),

		readyBtn:           RegisterObject(readyBtn),
		notReadyBtn:        RegisterObject(notReadyBtn),
		clearBoardBtn:      RegisterObject(clearBoardBtn),
		opponentReadyLabel: RegisterObject(opponentReadyLabel),

		opponentBoard:   RegisterObject(opponentBoard),
		playerTurnLabel: RegisterObject(playerTurnLabel),

		theEndLabel: RegisterObject(theEndLabel),

		objects: GlobalGameObjects.Objects(),
	}

	game.InitScenes()
	game.currentScene = game.scenes[SceneMenu]
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
	cp := data.NewPoint(float32(cx), float32(cy))

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
