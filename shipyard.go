package main

import (
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"

	"github.com/mymmrac/battleship/core"
	"github.com/mymmrac/battleship/ui"
)

var allowedShips = []int{
	0, // 1
	4, // 2
	3, // 3
	2, // 4
	1, // 5
}

// var allowedShips = []int{ // TODO: Remove
// 	2, // 1
// 	0, // 2
// 	0, // 3
// 	0, // 4
// 	0, // 5
// }

const shipyardBorder = false
const maxShipyardRowLen = 12 // 24

type Shipyard struct {
	core.BaseGameObject

	pos      core.Point[float32]
	board    *Board
	fontFace font.Face
	ships    []int
}

func NewShipyard(pos core.Point[float32], board *Board, fontFace font.Face) *Shipyard {
	return &Shipyard{
		BaseGameObject: core.NewBaseGameObject(),
		pos:            pos,
		board:          board,
		fontFace:       fontFace,
		ships:          make([]int, len(allowedShips)),
	}
}

func (s *Shipyard) Update(_ core.Point[float32]) {
	s.ships = s.shipsCount()
}

func (s *Shipyard) Draw(screen *ebiten.Image) {
	longestShip := len(allowedShips)

	// Border
	if shipyardBorder {
		size := longestShip*(longestShip+1)/2 + longestShip*2 - 1
		rows := float32(math.Ceil(float64(size)/maxShipyardRowLen))*2 - 1
		vector.StrokeRect(
			screen,
			s.pos.X-cellPaddingSize,
			s.pos.Y-cellPaddingSize,
			(maxShipyardRowLen)*cellSize+cellPaddingSize*2,
			(rows)*cellSize+cellPaddingSize*2,
			2,
			ui.BorderColor,
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
			pos.X,
			pos.Y,
			innerCellSize,
			innerCellSize,
			ui.MutedColor,
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
		ui.DrawCenteredText(
			screen,
			s.fontFace,
			strconv.Itoa(allowedShips[y]-s.ships[y]),
			int(pos.X+cellSize/2),
			int(pos.Y+cellSize/2),
			ui.TextDarkColor,
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
				pos.X,
				pos.Y,
				innerCellSize,
				innerCellSize,
				ui.ShipColor,
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

func (s *Shipyard) innerCellPos(x, y int) core.Point[float32] {
	return core.NewPoint(
		s.pos.X+float32(x)*cellSize+cellPaddingSize/2,
		s.pos.Y+float32(y)*cellSize+cellPaddingSize/2,
	)
}

func (s *Shipyard) cellPos(x, y int) core.Point[float32] {
	return core.NewPoint(
		s.pos.X+float32(x)*cellSize,
		s.pos.Y+float32(y)*cellSize,
	)
}

func (s *Shipyard) shipsCount() []int {
	ships := make([]int, len(allowedShips))
	visited := [cellsCount][cellsCount]bool{}

	for y := 0; y < cellsCount; y++ {
		for x := 0; x < cellsCount; x++ {
			if visited[y][x] {
				continue
			}
			visited[y][x] = true

			if s.board.At(x, y) == CellShip {
				l := 1

				dx := x + 1
				for dx < cellsCount && s.board.At(dx, y) == CellShip {
					visited[y][dx] = true
					l++
					dx++
				}

				if l == 1 {
					dy := y + 1
					for dy < cellsCount && s.board.At(x, dy) == CellShip {
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

func (s *Shipyard) ready() bool {
	for i, count := range s.ships {
		if allowedShips[i]-count != 0 {
			return false
		}
	}

	return true
}
