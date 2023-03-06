package main

import (
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SceneID int

const (
	sceneNone SceneID = iota
	sceneMenu
	sceneNewGame
	sceneJoinGame
	scenePlaceShips
	scenePlayerReady
)

type Scene struct {
	OnEnter  func()
	OnUpdate func()
	OnLeave  func()
}

func (g *Game) ChangeScene(id SceneID) {
	leave := g.currentScene.OnLeave
	if leave != nil {
		leave()
	}

	var ok bool
	g.currentScene, ok = g.scenes[id]
	if !ok {
		panic("unknown scene ID " + strconv.Itoa(int(id)))
	}

	enter := g.currentScene.OnEnter
	if enter != nil {
		enter()
	}
}

func (g *Game) InitScenes() {
	scenes := map[SceneID]*Scene{
		sceneMenu: {
			OnEnter: func() {
				g.newGameBtn.EnableAndShow()
				g.joinGameBtn.EnableAndShow()
				g.exitBtn.EnableAndShow()
			},
			OnUpdate: func() {
				if g.newGameBtn.clicked {
					g.ChangeScene(sceneNewGame)
					return
				}

				if g.joinGameBtn.clicked {
					g.ChangeScene(sceneJoinGame)
					return
				}

				if g.exitBtn.clicked {
					g.exit = true
				}
			},
			OnLeave: func() {
				g.newGameBtn.DisableAndHide()
				g.joinGameBtn.DisableAndHide()
				g.exitBtn.DisableAndHide()
			},
		},

		sceneNewGame: {
			OnEnter: func() {
				go func() {
					err := g.connector.StartNewGame()
					if err != nil {
						g.events <- NewEventError(EventNewGameStartFailed, err)
						return
					}

					g.events <- NewEventSignal(EventNewGameStarted)
				}()
			},
			OnUpdate: func() {
				event, ok := <-g.events
				if !ok {
					return
				}

				switch event.EventType() {
				case EventNewGameStarted:
					go func() {
						err := g.connector.WaitForConnection()
						if err != nil {
							g.events <- NewEventError(EventJoinGameFailed, err)
							return
						}

						g.events <- NewEventSignal(EventJoinedGame)
					}()

					// TODO: Make separate scene
					// g.ChangeScene(scenePlaceShips)
					// return

				case EventJoinedGame:
					g.ChangeScene(scenePlaceShips)
					return
				case EventNewGameStartFailed:
					errEvent := event.(EventError)
					fmt.Println(errEvent.err) // TODO: Fix me
					g.ChangeScene(sceneMenu)
					return
				default:
					panic("unexpected event type: " + strconv.Itoa(int(event.EventType())))
				}
			},
			OnLeave: nil,
		},

		sceneJoinGame: {
			OnEnter: func() {
				go func() {
					err := g.connector.JoinGame()
					if err != nil {
						g.events <- NewEventError(EventJoinGameFailed, err)
						return
					}

					g.events <- NewEventSignal(EventJoinedGame)
				}()
			},
			OnUpdate: func() {
				event, ok := <-g.events
				if !ok {
					return
				}

				switch event.EventType() {
				case EventJoinedGame:
					g.ChangeScene(scenePlaceShips)
					return
				case EventJoinGameFailed:
					errEvent := event.(EventError)
					fmt.Println(errEvent.err) // TODO: Fix me
					g.ChangeScene(sceneMenu)
					return
				default:
					panic("unexpected event type: " + strconv.Itoa(int(event.EventType())))
				}
			},
			OnLeave: nil,
		},

		scenePlaceShips: {
			OnEnter: func() {
				g.myBoard.EnableAndShow()
				g.myShipyard.EnableAndShow()
				g.clearBoardBtn.EnableAndShow()

				g.readyBtn.Disable()
				g.readyBtn.Show()
			},
			OnUpdate: func() {
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

				g.readyBtn.SetActive(g.myShipyard.ready())

				if g.readyBtn.clicked {
					g.ChangeScene(scenePlayerReady)
					return
				}
			},
			OnLeave: func() {
				g.myBoard.DisableAndHide()
				g.myShipyard.DisableAndHide()
				g.readyBtn.DisableAndHide()
				g.clearBoardBtn.DisableAndHide()
			},
		},

		scenePlayerReady: {
			OnEnter: func() {
				g.myBoard.Show()
				g.opponentBoard.EnableAndShow()
				g.notReadyBtn.EnableAndShow()
				g.clearBoardBtn.DisableAndHide()
			},
			OnUpdate: func() {
				if g.opponentBoard.hover && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					_ = g.opponentBoard.shoot(g.opponentBoard.hoverX, g.opponentBoard.hoverY)
				}

				if g.notReadyBtn.clicked {
					g.ChangeScene(scenePlaceShips)
				}
			},
			OnLeave: func() {
				g.myBoard.DisableAndHide()
				g.opponentBoard.DisableAndHide()
				g.notReadyBtn.DisableAndHide()
			},
		},
	}

	g.scenes = scenes
}
