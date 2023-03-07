package data

type Numeric interface {
	int | float32
}

type Point[T Numeric] struct {
	X T
	Y T
}

func NewPoint[T Numeric](x, y T) Point[T] {
	return Point[T]{
		X: x,
		Y: y,
	}
}

func (p Point[T]) Add(other Point[T]) Point[T] {
	return Point[T]{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

func (p Point[T]) Sub(other Point[T]) Point[T] {
	return Point[T]{
		X: p.X - other.X,
		Y: p.Y - other.Y,
	}
}
