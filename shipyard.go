package main

import (
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

// var allowedShips = []int{
// 	0, // 1
// 	4, // 2
// 	3, // 3
// 	2, // 4
// 	1, // 5
// }

var allowedShips = []int{ // TODO: Remove
	2, // 1
	0, // 2
	0, // 3
	0, // 4
	0, // 5
}

const shipyardBorder = false
const maxShipyardRowLen = 12 // 24

type shipyard struct {
	GameObject

	pos      point[float32]
	board    *board
	fontFace font.Face
	ships    []int
}

func newShipyard(pos point[float32], board *board, fontFace font.Face) *shipyard {
	return &shipyard{
		GameObject: NewGameObject(),
		pos:        pos,
		board:      board,
		fontFace:   fontFace,
		ships:      make([]int, len(allowedShips)),
	}
}

func (s *shipyard) Update(_ point[float32]) {
	s.ships = s.shipsCount()
}

func (s *shipyard) Draw(screen *ebiten.Image) {
	longestShip := len(allowedShips)

	// Border
	if shipyardBorder {
		size := longestShip*(longestShip+1)/2 + longestShip*2 - 1
		rows := float32(math.Ceil(float64(size)/maxShipyardRowLen))*2 - 1
		vector.StrokeRect(
			screen,
			s.pos.x-cellPaddingSize,
			s.pos.y-cellPaddingSize,
			(maxShipyardRowLen)*cellSize+cellPaddingSize*2,
			(rows)*cellSize+cellPaddingSize*2,
			2,
			borderColor,
		)
	}

	// Count cells
	row, col := 0, 0
	for y := 0; y < longestShip; y++ {
		if col+y+1+1 > maxShipyardRowLen {
			col = 0
			row += 2
		}

		pos := s.innerCellPos(col, row)
		vector.DrawFilledRect(
			screen,
			pos.x,
			pos.y,
			innerCellSize,
			innerCellSize,
			mutedColor,
		)

		col += y + 1 + 2
		if col > maxShipyardRowLen {
			col = 0
			row += 2
		}
	}

	// Count cells text
	row, col = 0, 0
	for y := 0; y < longestShip; y++ {
		if col+y+1+1 > maxShipyardRowLen {
			col = 0
			row += 2
		}

		pos := s.cellPos(col, row)
		DrawCenteredText(
			screen,
			s.fontFace,
			strconv.Itoa(allowedShips[y]-s.ships[y]),
			int(pos.x+cellSize/2),
			int(pos.y+cellSize/2),
			textDarkColor,
		)

		col += y + 1 + 2
		if col > maxShipyardRowLen {
			col = 0
			row += 2
		}
	}

	// Ship cells
	row, col = 0, 0
	for y := 0; y < longestShip; y++ {
		if col+y+1+1 > maxShipyardRowLen {
			col = 0
			row += 2
		}

		for x := 0; x < y+1; x++ {
			pos := s.innerCellPos(col+1, row)
			vector.DrawFilledRect(
				screen,
				pos.x,
				pos.y,
				innerCellSize,
				innerCellSize,
				shipColor,
			)

			col++
			if col > maxShipyardRowLen {
				col = 0
				row += 2
			}
		}

		col += 2
		if col > maxShipyardRowLen {
			col = 0
			row += 2
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

func (s *shipyard) shipsCount() []int {
	ships := make([]int, len(allowedShips))
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

func (s *shipyard) ready() bool {
	for i, count := range s.ships {
		if allowedShips[i]-count != 0 {
			return false
		}
	}

	return true
}
