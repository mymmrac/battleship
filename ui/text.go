package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

func DrawCenteredText(screen *ebiten.Image, font font.Face, s string, cx, cy int, clr color.Color) {
	if len(s) == 0 {
		return
	}

	bounds := text.BoundString(font, s)
	x, y := cx-bounds.Min.X-bounds.Dx()/2, cy-bounds.Min.Y-bounds.Dy()/2
	text.Draw(screen, s, font, x, y, clr)
}

func DrawTopLeftText(screen *ebiten.Image, font font.Face, s string, lx, ly int, clr color.Color) {
	if len(s) == 0 {
		return
	}

	bounds := text.BoundString(font, s)
	x, y := lx-bounds.Min.X, ly-bounds.Min.Y
	text.Draw(screen, s, font, x, y, clr)
}
