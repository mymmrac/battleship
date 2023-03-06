package main

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"

	"github.com/mymmrac/battleship/core"
	"github.com/mymmrac/battleship/ui"
)

const cellsCount = 10
const cellSize float32 = 32
const cellPaddingSize float32 = 4
const innerCellSize = cellSize - cellPaddingSize

type cellKind int

const (
	cellEmpty cellKind = iota
	cellShip
	cellMiss
	cellShipHit
)

type Board struct {
	core.BaseGameObject

	pos      core.Point[float32]
	cells    [cellsCount][cellsCount]cellKind
	fontFace font.Face

	hover  bool
	hoverX int
	hoverY int
}

func NewBoard(pos core.Point[float32], fontFace font.Face) *Board {
	return &Board{
		BaseGameObject: core.NewBaseGameObject(),
		pos:            pos,
		cells:          [cellsCount][cellsCount]cellKind{},
		fontFace:       fontFace,
	}
}

func (b *Board) Update(cp core.Point[float32]) {
	b.hoverX, b.hoverY, b.hover = b.cellOn(cp)
}

func (b *Board) CursorPointer() bool {
	return b.hover
}

func (b *Board) Disable() {
	b.hover = false
	b.BaseGameObject.Disable()
}

func (b *Board) Draw(screen *ebiten.Image) {
	// Border
	vector.StrokeRect(
		screen,
		b.pos.X-cellPaddingSize,
		b.pos.Y-cellPaddingSize,
		(cellsCount+1)*cellSize+cellPaddingSize*2,
		(cellsCount+1)*cellSize+cellPaddingSize*2,
		2,
		ui.BorderColor,
	)

	// Outer cells
	for y := 0; y < cellsCount+1; y++ {
		for x := 0; x < cellsCount+1; x++ {
			if x != 0 && y != 0 {
				continue
			}

			pos := b.innerCellPos(x, y)
			vector.DrawFilledRect(
				screen,
				pos.X,
				pos.Y,
				innerCellSize,
				innerCellSize,
				ui.MutedColor,
			)
		}
	}

	// Outer cells text
	for y := 0; y < cellsCount+1; y++ {
		for x := 0; x < cellsCount+1; x++ {
			if x != 0 && y != 0 {
				continue
			}

			var cellText string
			if x >= 1 && y == 0 {
				cellText = string(rune('A' + x - 1))
			} else if x == 0 && y >= 1 {
				cellText = strconv.Itoa(y)
			}

			pos := b.cellPos(x, y)
			ui.DrawCenteredText(
				screen,
				b.fontFace,
				cellText,
				int(pos.X+cellSize/2),
				int(pos.Y+cellSize/2),
				ui.TextDarkColor,
			)
		}
	}

	// Inner cells
	for y := 0; y < cellsCount; y++ {
		for x := 0; x < cellsCount; x++ {
			cell := b.cells[y][x]

			var clr color.Color
			switch cell {
			case cellEmpty:
				clr = ui.EmptyColor
			case cellShip:
				clr = ui.ShipColor
			case cellMiss:
				clr = ui.MissColor
			case cellShipHit:
				clr = ui.ShipHitColor
			default:
				panic("unreachable")
			}

			pos := b.innerCellPos(x+1, y+1)
			vector.DrawFilledRect(
				screen,
				pos.X,
				pos.Y,
				innerCellSize,
				innerCellSize,
				clr,
			)
		}
	}

	if b.hover {
		b.highlightCell(screen, b.hoverX, b.hoverY)
	}
}

func (b *Board) innerCellPos(x, y int) core.Point[float32] {
	return core.NewPoint(
		b.pos.X+float32(x)*cellSize+cellPaddingSize/2,
		b.pos.Y+float32(y)*cellSize+cellPaddingSize/2,
	)
}

func (b *Board) cellPos(x, y int) core.Point[float32] {
	return core.NewPoint(
		b.pos.X+float32(x)*cellSize,
		b.pos.Y+float32(y)*cellSize,
	)
}

func (b *Board) cellOn(p core.Point[float32]) (int, int, bool) {
	p = p.Sub(b.pos.Add(core.NewPoint(cellSize, cellSize)))

	for y := 0; y < cellsCount; y++ {
		for x := 0; x < cellsCount; x++ {
			cp := core.NewPoint(float32(x)*cellSize, float32(y)*cellSize)
			if cp.X <= p.X && p.X <= cp.X+cellSize &&
				cp.Y <= p.Y && p.Y <= cp.Y+cellSize {
				return x, y, true
			}
		}
	}

	return -1, -1, false
}

func (b *Board) shoot(x, y int) bool {
	switch b.cells[y][x] {
	case cellEmpty:
		b.cells[y][x] = cellMiss
		return true
	case cellShip:
		b.cells[y][x] = cellShipHit
		return true
	}

	return false
}

func (b *Board) placeShip(x, y int) {
	if b.cells[y][x] != cellEmpty {
		return
	}

	// Check diagonals
	if x > 0 { // LEFT
		if y > 0 { // UP
			if b.cells[y-1][x-1] == cellShip {
				return
			}
		}

		if y < cellsCount-1 { // DOWN
			if b.cells[y+1][x-1] == cellShip {
				return
			}
		}
	}
	if x < cellsCount-1 { // RIGHT
		if y > 0 { // UP
			if b.cells[y-1][x+1] == cellShip {
				return
			}
		}

		if y < cellsCount-1 { // DOWN
			if b.cells[y+1][x+1] == cellShip {
				return
			}
		}
	}

	// Check length
	length := 1
	for dx := x - 1; dx >= 0 && b.cells[y][dx] == cellShip; dx-- {
		length++
	}
	for dx := x + 1; dx < cellsCount && b.cells[y][dx] == cellShip; dx++ {
		length++
	}
	for dy := y - 1; dy >= 0 && b.cells[dy][x] == cellShip; dy-- {
		length++
	}
	for dy := y + 1; dy < cellsCount && b.cells[dy][x] == cellShip; dy++ {
		length++
	}
	if length > len(allowedShips) {
		return
	}

	b.cells[y][x] = cellShip
}

func (b *Board) removeShip(x, y int) {
	if b.cells[y][x] != cellShip {
		return
	}

	b.cells[y][x] = cellEmpty
}

func (b *Board) at(x, y int) cellKind {
	return b.cells[y][x]
}

func (b *Board) highlightCell(screen *ebiten.Image, x, y int) {
	pos := b.cellPos(x+1, y+1)
	vector.StrokeRect(screen, pos.X, pos.Y, cellSize, cellSize, 4, ui.HighlightColor)
}
