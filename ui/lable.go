package ui

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"

	"github.com/mymmrac/battleship/core"
)

type LabelAlignment int

const (
	LabelAlignmentCenter LabelAlignment = iota
	LabelAlignmentTopLeft
)

type Label struct {
	core.BaseGameObject

	pos       core.Point[float32]
	text      string
	fontFace  font.Face
	alignment LabelAlignment
}

func NewLabel(pos core.Point[float32], text string, fontFace font.Face) *Label {
	return &Label{
		BaseGameObject: core.NewBaseGameObject(),
		pos:            pos,
		text:           text,
		fontFace:       fontFace,
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
	default:
		panic("unknown alignment: " + strconv.Itoa(int(l.alignment)))
	}
}
