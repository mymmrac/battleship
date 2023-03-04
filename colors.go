package main

import "image/color"

var emptyColor = color.White

var mutedColor = color.Gray16{Y: 0xaaff}

var borderColor = color.White

var highlightColor = color.RGBA{
	R: 236,
	G: 168,
	B: 105,
	A: 255,
}

var shipColor = color.RGBA{
	R: 83,
	G: 127,
	B: 231,
	A: 255,
}

var missColor = color.RGBA{
	R: 60,
	G: 64,
	B: 72,
	A: 255,
}

var shipHitColor = color.RGBA{
	R: 245,
	G: 80,
	B: 80,
	A: 255,
}

var textColor = color.Black
