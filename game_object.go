package main

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameObject interface {
	Active() bool
	Update(cp point[float32])
	Visible() bool
	Draw(screen *ebiten.Image)
}

func RegisterObject[T GameObject](object T) T {
	GlobalGameObjects.Add(object)
	return object
}

var GlobalGameObjects = NewGlobalObjects[GameObject]()

type GlobalObjects[T any] struct {
	lock   *sync.Mutex
	locked bool

	objects []T
}

func NewGlobalObjects[T any]() *GlobalObjects[T] {
	return &GlobalObjects[T]{
		lock:    &sync.Mutex{},
		objects: nil,
	}
}

func (o *GlobalObjects[T]) Acquire() {
	if o.locked {
		panic("global objects already used")
	}

	o.lock.Lock()
	o.locked = true
	o.objects = make([]T, 0)
}

func (o *GlobalObjects[T]) Release() {
	if o.objects != nil {
		o.objects = nil
		o.locked = false
		o.lock.Unlock()
	}
}

func (o *GlobalObjects[T]) Add(object T) {
	if !o.locked {
		panic("global objects not acquired")
	}

	o.objects = append(o.objects, object)
}

func (o *GlobalObjects[T]) Objects() []T {
	if !o.locked {
		panic("global objects not acquired")
	}

	return o.objects
}
