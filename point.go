package main

type numeric interface {
	int | float32
}

type point[T numeric] struct {
	x T
	y T
}

func newPoint[T numeric](x, y T) point[T] {
	return point[T]{
		x: x,
		y: y,
	}
}

func (p point[T]) add(other point[T]) point[T] {
	return point[T]{
		x: p.x + other.x,
		y: p.y + other.y,
	}
}

func (p point[T]) sub(other point[T]) point[T] {
	return point[T]{
		x: p.x - other.x,
		y: p.y - other.y,
	}
}
