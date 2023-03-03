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

type cellKind int

const (
	cellEmpty cellKind = iota
	cellShip
	cellMiss
	cellShipHit
)

type board struct {
	pos      point
	cells    [10][10]cellKind
	fontFace font.Face
}

func newBoard(pos point, fontFace font.Face) *board {
	c := [cellsCount][cellsCount]cellKind{}

	c[1][1] = cellShip
	c[1][3] = cellMiss
	c[1][5] = cellShipHit

	return &board{
		pos:      pos,
		cells:    c,
		fontFace: fontFace,
	}
}

func (b *board) draw(screen *ebiten.Image) {
	// Border
	vector.StrokeRect(
		screen,
		b.pos.x-cellPaddingSize,
		b.pos.y-cellPaddingSize,
		(cellsCount+1)*cellSize+cellPaddingSize*2,
		(cellsCount+1)*cellSize+cellPaddingSize*2,
		2,
		color.White,
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
				b.innerCellSize(),
				b.innerCellSize(),
				color.Gray16{Y: 0xaaff},
			)

			var cellText string
			if x >= 1 && y == 0 {
				cellText = strconv.Itoa(x)
			} else if x == 0 && y >= 1 {
				cellText = string(rune('A' + y - 1))
			}

			DrawCenteredText(
				screen,
				b.fontFace,
				cellText,
				int(b.pos.x+float32(x)*cellSize+cellSize/2),
				int(b.pos.y+float32(y)*cellSize+cellSize/2),
				color.Black,
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
				clr = color.White
			case cellShip:
				clr = color.RGBA{
					R: 83,
					G: 127,
					B: 231,
					A: 255,
				}
			case cellMiss:
				clr = color.RGBA{
					R: 60,
					G: 64,
					B: 72,
					A: 255,
				}
			case cellShipHit:
				clr = color.RGBA{
					R: 245,
					G: 80,
					B: 80,
					A: 255,
				}
			default:
				panic("unreachable")
			}

			pos := b.innerCellPos(x+1, y+1)
			vector.DrawFilledRect(
				screen,
				pos.x,
				pos.y,
				b.innerCellSize(),
				b.innerCellSize(),
				clr,
			)
		}
	}
}

func (b *board) innerCellPos(x, y int) point {
	return newPoint(
		b.pos.x+float32(x)*cellSize+cellPaddingSize/2,
		b.pos.y+float32(y)*cellSize+cellPaddingSize/2,
	)
}

func (b *board) innerCellSize() float32 {
	return cellSize - cellPaddingSize
}

func (b *board) cellOn(p point) (int, int, bool) {
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
