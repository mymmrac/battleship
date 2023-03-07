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

	"github.com/mymmrac/battleship/events"
	"github.com/mymmrac/battleship/server/api"
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
	SceneTheEnd
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
					g.grpcConn, err = grpc.Dial(g.serverAddr+":"+g.serverPort,
						grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						g.events <- events.NewGameEventError(events.GameEventNewGameStartFailed, err)
						return
					}

					client := api.NewEventManagerClient(g.grpcConn)
					g.eventManager, err = NewEventManagerClient(client)
					if err != nil {
						g.events <- events.NewGameEventError(events.GameEventNewGameStartFailed, err)
						return
					}

					err = g.eventManager.NewGame()
					if err != nil {
						g.events <- events.NewGameEventError(events.GameEventNewGameStartFailed, err)
						return
					}

					time.Sleep(time.Second)
					g.events <- events.NewGameEventSignal(events.GameEventNewGameStarted)

					// TODO: Move to separate place
					err = g.eventManager.HandleGameEvents(g.events)
					if err != nil {
						panic(err)
					}
				}()
			},
			OnUpdate: func() {
				var event events.GameEvent
				select {
				case event = <-g.events:
				// Pass
				default:
					return
				}

				switch event.EventType() {
				case events.GameEventNewGameStarted:
					g.newGameLoadingLabel.SetText("Waiting for other player to join...")

					// TODO: Make separate scene
					// g.ChangeScene(sceneWaitForPlayer)
					// return
				case events.GameEventFromServer:
					serverEvent := event.(events.ServerEvent)

					var signalEvent events.GameEventSignal
					err := json.Unmarshal(serverEvent.Data, &signalEvent)
					if err != nil {
						fmt.Println(err) // TODO: Fix me
						return
					}

					if signalEvent.Type == events.GameEventJoinedGame {
						g.ChangeScene(ScenePlaceShips)
						return
					}
				case events.GameEventNewGameStartFailed:
					errEvent := event.(events.GameEventError)
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
					g.grpcConn, err = grpc.Dial(g.serverAddr+":"+g.serverPort,
						grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						g.events <- events.NewGameEventError(events.GameEventJoinGameFailed, err)
						return
					}

					client := api.NewEventManagerClient(g.grpcConn)
					g.eventManager, err = NewEventManagerClient(client)
					if err != nil {
						g.events <- events.NewGameEventError(events.GameEventJoinGameFailed, err)
						return
					}

					games, err := g.eventManager.ListGames()
					if err != nil {
						g.events <- events.NewGameEventError(events.GameEventJoinGameFailed, err)
						return
					}

					err = g.eventManager.JoinGame(games[0])
					if err != nil {
						g.events <- events.NewGameEventError(events.GameEventJoinGameFailed, err)
						return
					}

					g.events <- events.NewGameEventSignal(events.GameEventJoinedGame)

					// TODO: Move to separate place
					err = g.eventManager.HandleGameEvents(g.events)
					if err != nil {
						panic(err)
					}
				}()
			},
			OnUpdate: func() {
				var event events.GameEvent
				select {
				case event = <-g.events:
				// Pass
				default:
					return
				}

				switch event.EventType() {
				case events.GameEventJoinedGame:
					g.ChangeScene(ScenePlaceShips)
					return
				case events.GameEventJoinGameFailed:
					errEvent := event.(events.GameEventError)
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

				var event events.GameEvent
				select {
				case event = <-g.events:
				// Pass
				default:
					return
				}

				switch event.EventType() {
				case events.GameEventFromServer:
					serverEvent := event.(events.ServerEvent)

					var signalEvent events.GameEventSignal
					err := json.Unmarshal(serverEvent.Data, &signalEvent)
					if err != nil {
						fmt.Println(err) // TODO: Fix me
						return
					}

					if signalEvent.Type == events.GameEventPlayerReady {
						g.opponentReady = true
						g.opponentReadyLabel.SetText("Opponent: ready")
						return
					} else if signalEvent.Type == events.GameEventPlayerNotReady {
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
					err := g.eventManager.SendGameEvent(events.NewGameEventSignal(events.GameEventPlayerReady))
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

				var event events.GameEvent
				select {
				case event = <-g.events:
				// Pass
				default:
					return
				}

				switch event.EventType() {
				case events.GameEventFromServer:
					serverEvent := event.(events.ServerEvent)

					var signalEvent events.GameEventSignal
					err := json.Unmarshal(serverEvent.Data, &signalEvent)
					if err != nil {
						fmt.Println(err) // TODO: Fix me
						return
					}

					if signalEvent.Type == events.GameEventPlayerReady {
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

					err := g.eventManager.SendGameEvent(events.NewGameEventSignal(events.GameEventPlayerNotReady))
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
					err := g.eventManager.SendGameEvent(events.NewGameEventCoord(pos))
					if err != nil {
						fmt.Println(err) // TODO: Fix me
						return
					}

					g.myTurn = false
					// g.playerTurnLabel.SetText("Opponent's Turn")
					g.lastShootPos = pos
				}

				var event events.GameEvent
				select {
				case event = <-g.events:
				// Pass
				default:
					return
				}

				switch event.EventType() {
				case events.GameEventFromServer:
					serverEvent := event.(events.ServerEvent)

					var signalEvent events.GameEventSignal
					if err := json.Unmarshal(serverEvent.Data, &signalEvent); err != nil {
						fmt.Println(err) // TODO: Fix me
						return
					}

					switch signalEvent.EventType() {
					case events.GameEventShoot:
						var coordEvent events.GameEventCoord
						if err := json.Unmarshal(serverEvent.Data, &coordEvent); err != nil {
							fmt.Println(err) // TODO: Fix me
							return
						}

						hit := false

						var sendEvent events.GameEvent
						switch g.myBoard.AtPos(coordEvent.Pos) {
						case CellEmpty:
							sendEvent = events.NewGameEventSignal(events.GameEventMiss)
							g.myBoard.SetAt(coordEvent.Pos, CellMiss)
						case CellShip:
							hit = true

							sendEvent = events.NewGameEventSignal(events.GameEventHit)
							g.myBoard.SetAt(coordEvent.Pos, CellShipHit)

							if g.myBoard.FillIfDestroyed(coordEvent.Pos) {
								sendEvent = events.NewGameEventSignal(events.GameEventDestroyed)
							}
						}

						go func() {
							if err := g.eventManager.SendGameEvent(sendEvent); err != nil {
								fmt.Println(err) // TODO: Fix
								return
							}
						}()

						g.myTurn = !hit

						if !g.myBoard.HasAlive() {
							go func() {
								if err := g.eventManager.SendGameEvent(events.NewGameEventSignal(events.GameEventGameEnded)); err != nil {
									fmt.Println(err) // TODO: Fix
									return
								}
							}()

							g.ChangeScene(SceneTheEnd)
							return
						}
					case events.GameEventMiss:
						g.opponentBoard.SetAt(g.lastShootPos, CellMiss)
					case events.GameEventHit:
						g.opponentBoard.SetAt(g.lastShootPos, CellShipHit)
						g.myTurn = true
					case events.GameEventDestroyed:
						g.opponentBoard.SetAt(g.lastShootPos, CellShipHit)
						_ = g.opponentBoard.FillIfDestroyed(g.lastShootPos)
						g.myTurn = true
					case events.GameEventGameEnded:
						g.ChangeScene(SceneTheEnd)
						return
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

		SceneTheEnd: {
			OnEnter: func() {
				if g.myBoard.HasAlive() {
					g.theEndLabel.SetText("You Won!")
				} else {
					g.theEndLabel.SetText("You Lose!")
				}
				g.theEndLabel.Show()
				g.myBoard.Show()
				g.myBoard.Disable()
				g.opponentBoard.Show()
				g.opponentBoard.Disable()
			},
			OnUpdate: func() {},
			OnLeave: func() {
				g.theEndLabel.Hide()
				g.myBoard.Hide()
				g.opponentBoard.Hide()
			},
		},
	}

	g.scenes = scenes
}
