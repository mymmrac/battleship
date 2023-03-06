package ui

import "image/color"

var EmptyColor = color.White

var MutedColor = color.Gray16{Y: 0xaaff}

var BorderColor = color.White

var HighlightColor = color.RGBA{
	R: 236,
	G: 168,
	B: 105,
	A: 255,
}

var ShipColor = color.RGBA{
	R: 83,
	G: 127,
	B: 231,
	A: 255,
}

var MissColor = color.RGBA{
	R: 60,
	G: 64,
	B: 72,
	A: 255,
}

var ShipHitColor = color.RGBA{
	R: 245,
	G: 80,
	B: 80,
	A: 255,
}

var TextDarkColor = color.Black
var TextLightColor = color.White
