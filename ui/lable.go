package ui

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"

	"github.com/mymmrac/battleship/core"
	"github.com/mymmrac/battleship/data"
)

type LabelAlignment int

const (
	_ LabelAlignment = iota
	LabelAlignmentCenter
	LabelAlignmentTopLeft
	LabelAlignmentTopCenter
)

type Label struct {
	core.BaseGameObject

	pos       data.Point[float32]
	text      string
	fontFace  font.Face
	alignment LabelAlignment
}

func NewLabel(pos data.Point[float32], text string, fontFace font.Face) *Label {
	return &Label{
		BaseGameObject: core.NewBaseGameObject(),
		pos:            pos,
		text:           text,
		fontFace:       fontFace,
		alignment:      LabelAlignmentTopLeft,
	}
}

func (l *Label) SetAlignment(alignment LabelAlignment) {
	l.alignment = alignment
}

func (l *Label) SetText(text string) {
	l.text = text
}

func (l *Label) Draw(screen *ebiten.Image) {
	switch l.alignment {
	case LabelAlignmentCenter:
		DrawCenteredText(screen, l.fontFace, l.text, int(l.pos.X), int(l.pos.Y), TextLightColor)
	case LabelAlignmentTopLeft:
		DrawTopLeftText(screen, l.fontFace, l.text, int(l.pos.X), int(l.pos.Y), TextLightColor)
	case LabelAlignmentTopCenter:
		DrawTopCenterText(screen, l.fontFace, l.text, int(l.pos.X), int(l.pos.Y), TextLightColor)
	default:
		panic("unknown alignment: " + strconv.Itoa(int(l.alignment)))
	}
}
