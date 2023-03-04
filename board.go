package main

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
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

type board struct {
	pos      point[float32]
	cells    [10][10]cellKind
	fontFace font.Face

	hover  bool
	hoverX int
	hoverY int
}

func newBoard(pos point[float32], fontFace font.Face) *board {
	return &board{
		pos:      pos,
		cells:    [cellsCount][cellsCount]cellKind{},
		fontFace: fontFace,
	}
}

func (b *board) update(cp point[float32]) {
	b.hoverX, b.hoverY, b.hover = b.cellOn(cp)
}

func (b *board) Draw(screen *ebiten.Image) {
	// Border
	vector.StrokeRect(
		screen,
		b.pos.x-cellPaddingSize,
		b.pos.y-cellPaddingSize,
		(cellsCount+1)*cellSize+cellPaddingSize*2,
		(cellsCount+1)*cellSize+cellPaddingSize*2,
		2,
		borderColor,
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
				pos.x,
				pos.y,
				innerCellSize,
				innerCellSize,
				mutedColor,
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
			DrawCenteredText(
				screen,
				b.fontFace,
				cellText,
				int(pos.x+cellSize/2),
				int(pos.y+cellSize/2),
				textDarkColor,
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
				clr = emptyColor
			case cellShip:
				clr = shipColor
			case cellMiss:
				clr = missColor
			case cellShipHit:
				clr = shipHitColor
			default:
				panic("unreachable")
			}

			pos := b.innerCellPos(x+1, y+1)
			vector.DrawFilledRect(
				screen,
				pos.x,
				pos.y,
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

func (b *board) innerCellPos(x, y int) point[float32] {
	return newPoint(
		b.pos.x+float32(x)*cellSize+cellPaddingSize/2,
		b.pos.y+float32(y)*cellSize+cellPaddingSize/2,
	)
}

func (b *board) cellPos(x, y int) point[float32] {
	return newPoint(
		b.pos.x+float32(x)*cellSize,
		b.pos.y+float32(y)*cellSize,
	)
}

func (b *board) cellOn(p point[float32]) (int, int, bool) {
	p = p.sub(b.pos.add(newPoint(cellSize, cellSize)))

	for y := 0; y < cellsCount; y++ {
		for x := 0; x < cellsCount; x++ {
			cp := newPoint(float32(x)*cellSize, float32(y)*cellSize)
			if cp.x <= p.x && p.x <= cp.x+cellSize &&
				cp.y <= p.y && p.y <= cp.y+cellSize {
				return x, y, true
			}
		}
	}

	return -1, -1, false
}

func (b *board) shoot(x, y int) bool {
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

func (b *board) placeShip(x, y int) {
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

func (b *board) removeShip(x, y int) {
	if b.cells[y][x] != cellShip {
		return
	}

	b.cells[y][x] = cellEmpty
}

func (b *board) at(x, y int) cellKind {
	return b.cells[y][x]
}

func (b *board) highlightCell(screen *ebiten.Image, x, y int) {
	pos := b.cellPos(x+1, y+1)
	vector.StrokeRect(screen, pos.x, pos.y, cellSize, cellSize, 4, highlightColor)
}
