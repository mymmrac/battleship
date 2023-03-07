package core

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/mymmrac/battleship/data"
)

type BaseGameObject struct {
	active  bool
	visible bool
}

func NewBaseGameObject() BaseGameObject {
	return BaseGameObject{
		active:  false,
		visible: false,
	}
}

func (o *BaseGameObject) Update(_ data.Point[float32]) {}

func (o *BaseGameObject) Draw(_ *ebiten.Image) {}

func (o *BaseGameObject) Active() bool {
	return o.active
}

func (o *BaseGameObject) SetActive(is bool) {
	o.active = is
}

func (o *BaseGameObject) Enable() {
	o.active = true
}

func (o *BaseGameObject) Disable() {
	o.active = false
}

func (o *BaseGameObject) Visible() bool {
	return o.visible
}

func (o *BaseGameObject) SetVisible(is bool) {
	o.visible = is
}

func (o *BaseGameObject) Show() {
	o.visible = true
}

func (o *BaseGameObject) Hide() {
	o.visible = false
}

func (o *BaseGameObject) CursorPointer() bool {
	return false
}

func (o *BaseGameObject) EnableAndShow() {
	o.active = true
	o.visible = true
}

func (o *BaseGameObject) DisableAndHide() {
	o.active = false
	o.visible = false
}
