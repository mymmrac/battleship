package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

const buttonPadding float32 = 4

type button struct {
	BaseGameObject

	pos      point[float32]
	width    float32
	height   float32
	text     string
	fontFace font.Face

	hover   bool
	clicked bool
}

func newButton(pos point[float32], width float32, height float32, text string, fontFace font.Face) *button {
	return &button{
		BaseGameObject: NewBaseGameObject(),
		pos:            pos,
		width:          width,
		height:         height,
		text:           text,
		fontFace:       fontFace,
	}
}

func (b *button) Update(cp point[float32]) {
	b.hover = b.pos.x <= cp.x && cp.x <= b.pos.x+b.width &&
		b.pos.y <= cp.y && cp.y <= b.pos.y+b.height

	b.clicked = b.hover && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}

func (b *button) Disable() {
	b.hover = false
	b.clicked = false
	b.BaseGameObject.Disable()
}

func (b *button) Draw(screen *ebiten.Image) {
	// Border
	clr := borderColor
	if !b.Active() {
		clr = mutedColor
	}
	vector.StrokeRect(
		screen,
		b.pos.x,
		b.pos.y,
		b.width,
		b.height,
		2,
		clr,
	)

	// Background
	if b.hover {
		vector.DrawFilledRect(
			screen,
			b.pos.x+buttonPadding,
			b.pos.y+buttonPadding,
			b.width-buttonPadding*2,
			b.height-buttonPadding*2,
			mutedColor,
		)
	}

	// Text
	clr = textLightColor
	if b.hover {
		clr = textDarkColor
	}
	if !b.Active() {
		clr = mutedColor
	}
	DrawCenteredText(screen, b.fontFace, b.text, int(b.pos.x+b.width/2), int(b.pos.y+b.height/2), clr)
}
