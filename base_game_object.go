package main

import "github.com/hajimehoshi/ebiten/v2"

type BaseGameObject struct {
	active  bool
	visible bool
}

func NewBaseGameObject() BaseGameObject {
	return BaseGameObject{
		active:  true,
		visible: true,
	}
}

func (o *BaseGameObject) Update(_ point[float32]) {}

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
