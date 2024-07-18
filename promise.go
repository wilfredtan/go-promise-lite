package promise

import (
	"fmt"
	"sync"
)

type Promise[T any] struct {
	value T
	err   error
	once  sync.Once
	ch    chan any
}

// Waits for the Promise to resolve or reject and returns the resolved value or an error.
func (p *Promise[T]) Await() (T, error) {
	<-p.ch
	return p.value, p.err
}

// Resolves a promise by updating it with a value.
func (p *Promise[T]) resolve(value T) {
	p.once.Do(func() {
		p.value = value
		close(p.ch)
	})
}

// Rejects a promise by updating it with an error.
func (p *Promise[T]) reject(err error) {
	p.once.Do(func() {
		p.err = err
		close(p.ch)
	})
}

// Handles panics that occur during the execution of the Promise.
func (p *Promise[T]) panic_handler() {
	var unknown_exception any
	if unknown_exception = recover(); unknown_exception == nil {
		return
	}
	switch typed_exception := unknown_exception.(type) {
	case error:
		p.reject(typed_exception)
	default:
		p.reject(fmt.Errorf("unhandled error: %v", typed_exception))
	}
}

// New creates a new Promise with the provided resolver function.
func New[R any](resolver func(func(R), func(error))) *Promise[R] {
	promise := Promise[R]{
		ch:   make(chan any, 1),
		once: sync.Once{},
	}
	go func() {
		defer promise.panic_handler()
		resolver(promise.resolve, promise.reject)
	}()
	return &promise
}
