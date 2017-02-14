package pool

import "sync"

func NewPool(f func() interface{}) *sync.Pool {
	return &sync.Pool{
		New: f,
	}
}
