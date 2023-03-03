package main

type point struct {
	x float32
	y float32
}

func newPoint(x, y float32) point {
	return point{
		x: x,
		y: y,
	}
}

func (p point) add(other point) point {
	return point{
		x: p.x + other.x,
		y: p.y + other.y,
	}
}

func (p point) sub(other point) point {
	return point{
		x: p.x - other.x,
		y: p.y - other.y,
	}
}
