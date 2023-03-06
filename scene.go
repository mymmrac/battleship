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
						g.myBoard.placeShip(g.myBoard.hoverPos)
					}

					if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
						g.myBoard.removeShip(g.myBoard.hoverPos)
					}
				}

				if g.clearBoardBtn.Clicked() {
					g.myBoard.cells = [cellsCount][cellsCount]CellKind{}
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
					g.playerTurnLabel.SetText("Opponent's Turn")
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
						g.myTurn = true
						g.playerTurnLabel.SetText("Your Turn")
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
				g.playerTurnLabel.Show()
			},
			OnUpdate: func() {
				pos := g.opponentBoard.hoverPos
				if g.myTurn && g.opponentBoard.hover && g.opponentBoard.canShoot(pos) &&
					inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					err := g.eventManager.SendGameEvent(NewGameEventCoord(pos))
					if err != nil {
						fmt.Println(err) // TODO: Fix me
						return
					}

					g.myTurn = false
					// g.playerTurnLabel.SetText("Opponent's Turn")
					g.lastShootPos = pos
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
					if err := json.Unmarshal(serverEvent.Data, &signalEvent); err != nil {
						fmt.Println(err) // TODO: Fix me
						return
					}

					switch signalEvent.EventType() {
					case GameEventShoot:
						var coordEvent GameEventCoord
						if err := json.Unmarshal(serverEvent.Data, &coordEvent); err != nil {
							fmt.Println(err) // TODO: Fix me
							return
						}

						hit := false

						var sendEvent GameEvent
						switch g.myBoard.AtPos(coordEvent.Pos) {
						case CellEmpty:
							sendEvent = NewGameEventSignal(GameEventMiss)
							g.myBoard.SetAt(coordEvent.Pos, CellMiss)
						case CellShip:
							hit = true

							sendEvent = NewGameEventSignal(GameEventHit)
							g.myBoard.SetAt(coordEvent.Pos, CellShipHit)

							if g.myBoard.FillIfDestroyed(coordEvent.Pos) {
								sendEvent = NewGameEventSignal(GameEventDestroyed)
							}
						}

						go func() {
							if err := g.eventManager.SendGameEvent(sendEvent); err != nil {
								fmt.Println(err) // TODO: Fix
								return
							}
						}()

						g.myTurn = !hit
					case GameEventMiss:
						g.opponentBoard.SetAt(g.lastShootPos, CellMiss)
					case GameEventHit:
						g.opponentBoard.SetAt(g.lastShootPos, CellShipHit)
						g.myTurn = true
					case GameEventDestroyed:
						g.opponentBoard.SetAt(g.lastShootPos, CellShipHit)
						_ = g.opponentBoard.FillIfDestroyed(g.lastShootPos)
						g.myTurn = true
						// TODO: Mark all empty spots
					}
				default:
					panic("unexpected event type: " + strconv.Itoa(int(event.EventType())))
				}

				if g.myTurn {
					g.playerTurnLabel.SetText("Your Turn")
				} else {
					g.playerTurnLabel.SetText("Opponent's Turn")
				}
			},
			OnLeave: func() {
				g.myBoard.DisableAndHide()
				g.opponentBoard.DisableAndHide()
				g.playerTurnLabel.Hide()
			},
		},
	}

	g.scenes = scenes
}
