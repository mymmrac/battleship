package main

import "github.com/hajimehoshi/ebiten/v2"

const longestShip = 5

var allowedShips = [longestShip]int{
	0, // 1
	4, // 2
	3, // 3
	2, // 4
	1, // 5
}

type shipyard struct {
	pos   point[float32]
	board *board
}

func newShipyard(pos point[float32], board *board) *shipyard {
	return &shipyard{
		pos:   pos,
		board: board,
	}
}

func (s *shipyard) draw(screen *ebiten.Image) {

}

func (s *shipyard) countShips() [longestShip]int {
	ships := [longestShip]int{}

	visited := [cellsCount][cellsCount]bool{}

	for y := 0; y < cellsCount; y++ {
		for x := 0; x < cellsCount; x++ {
			if visited[y][x] {
				continue
			}
			visited[y][x] = true

			if s.board.at(x, y) == cellShip {
				l := 1

				dx := x + 1
				for dx < cellsCount && s.board.at(dx, y) == cellShip {
					visited[y][dx] = true
					l++
					dx++
				}

				if l == 1 {
					dy := y + 1
					for dy < cellsCount && s.board.at(x, dy) == cellShip {
						visited[dy][x] = true
						l++
						dy++
					}
				}

				ships[l-1]++
			}
		}
	}

	return ships
}
