package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/mymmrac/battleship/api"
)

type SceneID int

const (
	_ SceneID = iota
	SceneMenu
	SceneNewGame
	SceneJoinGame
	ScenePlaceShips
	ScenePlayerReady
	SceneTheGame
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
		SceneMenu: {
			OnEnter: func() {
				g.newGameBtn.EnableAndShow()
				g.joinGameBtn.EnableAndShow()
				g.exitBtn.EnableAndShow()
			},
			OnUpdate: func() {
				if g.newGameBtn.Clicked() {
					g.ChangeScene(SceneNewGame)
					return
				}

				if g.joinGameBtn.Clicked() {
					g.ChangeScene(SceneJoinGame)
					return
				}

				if g.exitBtn.Clicked() {
					g.exit = true
				}
			},
			OnLeave: func() {
				g.newGameBtn.DisableAndHide()
				g.joinGameBtn.DisableAndHide()
				g.exitBtn.DisableAndHide()
			},
		},

		SceneNewGame: {
			OnEnter: func() {
				g.newGameLoadingLabel.Show()
				g.newGameLoadingLabel.SetText("Creating new game...")

				go func() {
					var err error
					// TODO: Close connection
					g.grpcConn, err = grpc.Dial(grpcAddr+":"+grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						g.events <- NewGameEventError(GameEventNewGameStartFailed, err)
						return
					}

					client := api.NewEventManagerClient(g.grpcConn)
					g.eventManager, err = NewEventManagerClient(client)
					if err != nil {
						g.events <- NewGameEventError(GameEventNewGameStartFailed, err)
						return
					}

					err = g.eventManager.NewGame()
					if err != nil {
						g.events <- NewGameEventError(GameEventNewGameStartFailed, err)
						return
					}

					time.Sleep(time.Second)
					g.events <- NewGameEventSignal(GameEventNewGameStarted)

					// TODO: Move to separate place
					err = g.eventManager.HandleGameEvents(g.events)
					if err != nil {
						panic(err)
					}
				}()
			},
			OnUpdate: func() {
				var event GameEvent
				select {
				case event = <-g.events:
				// Pass
				default:
					return
				}

				switch event.EventType() {
				case GameEventNewGameStarted:
					g.newGameLoadingLabel.SetText("Waiting for other player to join...")

					// TODO: Make separate scene
					// g.ChangeScene(sceneWaitForPlayer)
					// return
				case GameEventFromServer:
					serverEvent := event.(ServerEvent)

					var signalEvent GameEventSignal
					err := json.Unmarshal(serverEvent.Data, &signalEvent)
					if err != nil {
						fmt.Println(err) // TODO: Fix me
						return
					}

					if signalEvent.Type == GameEventJoinedGame {
						g.ChangeScene(ScenePlaceShips)
						return
					}
				case GameEventNewGameStartFailed:
					errEvent := event.(GameEventError)
					fmt.Println(errEvent.Err) // TODO: Fix me
					g.ChangeScene(SceneMenu)
					return
				default:
					panic("unexpected event type: " + strconv.Itoa(int(event.EventType())))
				}
			},
			OnLeave: func() {
				g.newGameLoadingLabel.Hide()
			},
		},

		SceneJoinGame: {
			OnEnter: func() {
				go func() {
					var err error
					// TODO: Close connection
					g.grpcConn, err = grpc.Dial(grpcAddr+":"+grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						g.events <- NewGameEventError(GameEventJoinGameFailed, err)
						return
					}

					client := api.NewEventManagerClient(g.grpcConn)
					g.eventManager, err = NewEventManagerClient(client)
					if err != nil {
						g.events <- NewGameEventError(GameEventJoinGameFailed, err)
						return
					}

					games, err := g.eventManager.ListGames()
					if err != nil {
						g.events <- NewGameEventError(GameEventJoinGameFailed, err)
						return
					}

					err = g.eventManager.JoinGame(games[0])
					if err != nil {
						g.events <- NewGameEventError(GameEventJoinGameFailed, err)
						return
					}

					g.events <- NewGameEventSignal(GameEventJoinedGame)

					// TODO: Move to separate place
					err = g.eventManager.HandleGameEvents(g.events)
					if err != nil {
						panic(err)
					}
				}()
			},
			OnUpdate: func() {
				var event GameEvent
				select {
				case event = <-g.events:
				// Pass
				default:
					return
				}

				switch event.EventType() {
				case GameEventJoinedGame:
					g.ChangeScene(ScenePlaceShips)
					return
				case GameEventJoinGameFailed:
					errEvent := event.(GameEventError)
					fmt.Println(errEvent.Err) // TODO: Fix me
					g.ChangeScene(SceneMenu)
					return
				default:
					panic("unexpected event type: " + strconv.Itoa(int(event.EventType())))
				}
			},
			OnLeave: nil,
		},

		ScenePlaceShips: {
			OnEnter: func() {
				g.myBoard.EnableAndShow()
				g.myShipyard.EnableAndShow()
				g.clearBoardBtn.EnableAndShow()
				g.opponentReadyLabel.Show()

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

				if g.clearBoardBtn.Clicked() {
					g.myBoard.cells = [cellsCount][cellsCount]cellKind{}
				}

				g.readyBtn.SetActive(g.myShipyard.ready())

				if g.readyBtn.Clicked() {
					g.ChangeScene(ScenePlayerReady)
					return
				}

				var event GameEvent
				select {
				case event = <-g.events:
				// Pass
				default:
					return
				}

				switch event.EventType() {
				case GameEventFromServer:
					serverEvent := event.(ServerEvent)

					var signalEvent GameEventSignal
					err := json.Unmarshal(serverEvent.Data, &signalEvent)
					if err != nil {
						fmt.Println(err) // TODO: Fix me
						return
					}

					if signalEvent.Type == GameEventPlayerReady {
						g.opponentReady = true
						g.opponentReadyLabel.SetText("Opponent: ready")
						return
					} else if signalEvent.Type == GameEventPlayerNotReady {
						g.opponentReady = false
						g.opponentReadyLabel.SetText("Opponent: not ready")
						return
					}
				default:
					panic("unexpected event type: " + strconv.Itoa(int(event.EventType())))
				}
			},
			OnLeave: func() {
				g.myBoard.DisableAndHide()
				g.myShipyard.DisableAndHide()
				g.readyBtn.DisableAndHide()
				g.clearBoardBtn.DisableAndHide()
				g.opponentReadyLabel.Hide()
			},
		},

		ScenePlayerReady: {
			OnEnter: func() {
				go func() {
					err := g.eventManager.SendGameEvent(NewGameEventSignal(GameEventPlayerReady))
					if err != nil {
						fmt.Println(err) // TODO: Fix me
						return
					}
				}()

				g.myBoard.Show()
				g.notReadyBtn.EnableAndShow()
				g.clearBoardBtn.DisableAndHide()
				g.opponentReadyLabel.Show()
			},
			OnUpdate: func() {
				if g.opponentReady {
					g.ChangeScene(SceneTheGame)
					return
				}

				if g.notReadyBtn.Clicked() {
					g.ChangeScene(ScenePlaceShips)
				}

				var event GameEvent
				select {
				case event = <-g.events:
				// Pass
				default:
					return
				}

				switch event.EventType() {
				case GameEventFromServer:
					serverEvent := event.(ServerEvent)

					var signalEvent GameEventSignal
					err := json.Unmarshal(serverEvent.Data, &signalEvent)
					if err != nil {
						fmt.Println(err) // TODO: Fix me
						return
					}

					if signalEvent.Type == GameEventPlayerReady {
						g.opponentReady = true
						g.ChangeScene(SceneTheGame)
						return
					}
				default:
					panic("unexpected event type: " + strconv.Itoa(int(event.EventType())))
				}
			},
			OnLeave: func() {
				g.myBoard.DisableAndHide()
				g.notReadyBtn.DisableAndHide()
				g.opponentReadyLabel.Hide()

				go func() {
					if g.opponentReady {
						return
					}

					err := g.eventManager.SendGameEvent(NewGameEventSignal(GameEventPlayerNotReady))
					if err != nil {
						fmt.Println(err) // TODO: Fix me
						return
					}
				}()
			},
		},

		SceneTheGame: {
			OnEnter: func() {
				g.myBoard.Show()
				g.opponentBoard.EnableAndShow()
			},
			OnUpdate: func() {
				if g.opponentBoard.hover && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					_ = g.opponentBoard.shoot(g.opponentBoard.hoverX, g.opponentBoard.hoverY)
				}
			},
			OnLeave: func() {
				g.myBoard.DisableAndHide()
				g.opponentBoard.DisableAndHide()
			},
		},
	}

	g.scenes = scenes
}
