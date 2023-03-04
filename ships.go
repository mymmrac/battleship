package main

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

const longestShip = 5

var allowedShips = [longestShip]int{
	0, // 1
	4, // 2
	3, // 3
	2, // 4
	1, // 5
}

type shipyard struct {
	pos      point[float32]
	board    *board
	fontFace font.Face
	counts   [longestShip]int
}

func newShipyard(pos point[float32], board *board, fontFace font.Face) *shipyard {
	return &shipyard{
		pos:      pos,
		board:    board,
		fontFace: fontFace,
		counts:   [longestShip]int{},
	}
}

func (s *shipyard) draw(screen *ebiten.Image) {
	// Border
	vector.StrokeRect(
		screen,
		s.pos.x-cellPaddingSize,
		s.pos.y-cellPaddingSize,
		(longestShip+1)*cellSize+cellPaddingSize*2,
		(longestShip)*cellSize+cellPaddingSize*2,
		2,
		borderColor,
	)

	// Count cells
	for y := 0; y < longestShip; y++ {
		pos := s.innerCellPos(0, y)
		vector.DrawFilledRect(
			screen,
			pos.x,
			pos.y,
			innerCellSize,
			innerCellSize,
			mutedColor,
		)
	}

	// Count cells text
	for y := 0; y < longestShip; y++ {
		pos := s.cellPos(0, y)
		DrawCenteredText(
			screen,
			s.fontFace,
			strconv.Itoa(allowedShips[y]-s.counts[y]),
			int(pos.x+cellSize/2),
			int(pos.y+cellSize/2),
			textColor,
		)
	}

	// Ship cells
	for y := 0; y < longestShip; y++ {
		for x := 0; x < y+1; x++ {
			pos := s.innerCellPos(x+1, y)
			vector.DrawFilledRect(
				screen,
				pos.x,
				pos.y,
				innerCellSize,
				innerCellSize,
				shipColor,
			)
		}
	}
}

func (s *shipyard) innerCellPos(x, y int) point[float32] {
	return newPoint(
		s.pos.x+float32(x)*cellSize+cellPaddingSize/2,
		s.pos.y+float32(y)*cellSize+cellPaddingSize/2,
	)
}

func (s *shipyard) cellPos(x, y int) point[float32] {
	return newPoint(
		s.pos.x+float32(x)*cellSize,
		s.pos.y+float32(y)*cellSize,
	)
}

func (s *shipyard) update() {
	s.counts = s.countShips()
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
