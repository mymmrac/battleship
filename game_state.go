package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameState int

const (
	statePlaceShips GameState = iota
	statePlayerReady
)

type Scene struct {
	OnEnter  func()
	OnUpdate func()
	OnLeave  func()
}

func (g *Game) ChangeScene(state GameState) {
	leave := g.currentScene.OnLeave
	if leave != nil {
		leave()
	}

	g.currentScene = g.scenes[state]

	enter := g.currentScene.OnEnter
	if enter != nil {
		enter()
	}
}

func (g *Game) InitScenes() {
	scenes := map[GameState]*Scene{
		statePlaceShips: {
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
					g.ChangeScene(statePlayerReady)
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
		statePlayerReady: {
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
					g.ChangeScene(statePlaceShips)
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
