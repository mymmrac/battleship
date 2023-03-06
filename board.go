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

type CellKind int

const (
	CellEmpty CellKind = iota
	CellShip
	CellMiss
	CellShipHit
)

type Board struct {
	core.BaseGameObject

	pos      core.Point[float32]
	cells    [cellsCount][cellsCount]CellKind
	fontFace font.Face

	hover    bool
	hoverPos core.Point[int]
}

func NewBoard(pos core.Point[float32], fontFace font.Face) *Board {
	return &Board{
		BaseGameObject: core.NewBaseGameObject(),
		pos:            pos,
		cells:          [cellsCount][cellsCount]CellKind{},
		fontFace:       fontFace,
	}
}

func (b *Board) Update(cp core.Point[float32]) {
	b.hoverPos.X, b.hoverPos.Y, b.hover = b.cellOn(cp)
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
			case CellEmpty:
				clr = ui.EmptyColor
			case CellShip:
				clr = ui.ShipColor
			case CellMiss:
				clr = ui.MissColor
			case CellShipHit:
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
		b.highlightCell(screen, b.hoverPos)
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

func (b *Board) canShoot(pos core.Point[int]) bool {
	return b.cells[pos.Y][pos.X] == CellEmpty
}

func (b *Board) placeShip(pos core.Point[int]) {
	x, y := pos.X, pos.Y

	if b.cells[y][x] != CellEmpty {
		return
	}

	// Check diagonals
	if x > 0 { // LEFT
		if y > 0 { // UP
			if b.cells[y-1][x-1] == CellShip {
				return
			}
		}

		if y < cellsCount-1 { // DOWN
			if b.cells[y+1][x-1] == CellShip {
				return
			}
		}
	}
	if x < cellsCount-1 { // RIGHT
		if y > 0 { // UP
			if b.cells[y-1][x+1] == CellShip {
				return
			}
		}

		if y < cellsCount-1 { // DOWN
			if b.cells[y+1][x+1] == CellShip {
				return
			}
		}
	}

	// Check length
	length := 1
	for dx := x - 1; dx >= 0 && b.cells[y][dx] == CellShip; dx-- {
		length++
	}
	for dx := x + 1; dx < cellsCount && b.cells[y][dx] == CellShip; dx++ {
		length++
	}
	for dy := y - 1; dy >= 0 && b.cells[dy][x] == CellShip; dy-- {
		length++
	}
	for dy := y + 1; dy < cellsCount && b.cells[dy][x] == CellShip; dy++ {
		length++
	}
	if length > len(allowedShips) {
		return
	}

	b.cells[y][x] = CellShip
}

func (b *Board) removeShip(pos core.Point[int]) {
	if b.cells[pos.Y][pos.X] != CellShip {
		return
	}

	b.cells[pos.Y][pos.X] = CellEmpty
}

func (b *Board) At(x, y int) CellKind {
	return b.cells[y][x]
}

func (b *Board) AtPos(pos core.Point[int]) CellKind {
	return b.cells[pos.Y][pos.X]
}

func (b *Board) highlightCell(screen *ebiten.Image, boardPos core.Point[int]) {
	pos := b.cellPos(boardPos.X+1, boardPos.Y+1)
	vector.StrokeRect(screen, pos.X, pos.Y, cellSize, cellSize, 4, ui.HighlightColor)
}

func (b *Board) SetAt(pos core.Point[int], kind CellKind) {
	b.cells[pos.Y][pos.X] = kind
}
