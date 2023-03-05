package main

type GameObject struct {
	active  bool
	visible bool
}

func NewGameObject() GameObject {
	return GameObject{
		active:  true,
		visible: true,
	}
}

func (o *GameObject) Active() bool {
	return o.active
}

func (o *GameObject) SetActive(is bool) {
	o.active = is
}

func (o *GameObject) Enable() {
	o.active = true
}

func (o *GameObject) Disable() {
	o.active = false
}

func (o *GameObject) Visible() bool {
	return o.visible
}

func (o *GameObject) SetVisible(is bool) {
	o.visible = is
}

func (o *GameObject) Show() {
	o.visible = true
}

func (o *GameObject) Hide() {
	o.visible = false
}
