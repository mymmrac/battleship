package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"

	"github.com/mymmrac/battleship/core"
)

const buttonPadding float32 = 4

type Button struct {
	core.BaseGameObject

	pos      core.Point[float32]
	width    float32
	height   float32
	text     string
	fontFace font.Face

	hover   bool
	clicked bool
}

func NewButton(pos core.Point[float32], width float32, height float32, text string, fontFace font.Face) *Button {
	return &Button{
		BaseGameObject: core.NewBaseGameObject(),
		pos:            pos,
		width:          width,
		height:         height,
		text:           text,
		fontFace:       fontFace,
	}
}

func (b *Button) Update(cp core.Point[float32]) {
	b.hover = b.pos.X <= cp.X && cp.X <= b.pos.X+b.width &&
		b.pos.Y <= cp.Y && cp.Y <= b.pos.Y+b.height

	b.clicked = b.hover && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}

func (b *Button) Disable() {
	b.hover = false
	b.clicked = false
	b.BaseGameObject.Disable()
}

func (b *Button) CursorPointer() bool {
	return b.hover
}

func (b *Button) Clicked() bool {
	return b.clicked
}

func (b *Button) Draw(screen *ebiten.Image) {
	// Border
	clr := BorderColor
	if !b.Active() {
		clr = MutedColor
	}
	vector.StrokeRect(
		screen,
		b.pos.X,
		b.pos.Y,
		b.width,
		b.height,
		2,
		clr,
	)

	// Background
	if b.hover {
		vector.DrawFilledRect(
			screen,
			b.pos.X+buttonPadding,
			b.pos.Y+buttonPadding,
			b.width-buttonPadding*2,
			b.height-buttonPadding*2,
			MutedColor,
		)
	}

	// Text
	clr = TextLightColor
	if b.hover {
		clr = TextDarkColor
	}
	if !b.Active() {
		clr = MutedColor
	}

	DrawCenteredText(screen, b.fontFace, b.text, int(b.pos.X+b.width/2), int(b.pos.Y+b.height/2), clr)
}
