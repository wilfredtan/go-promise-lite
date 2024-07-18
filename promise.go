package promise

import (
	"fmt"
)

type Promise[T any] struct {
	ch_resolved_value chan T
	ch_error          chan error
	resolved          bool
}

// Waits for the Promise to resolve or reject and returns the resolved value or an error.
func (p *Promise[T]) Await() (T, error) {
	var data T
	var err error
	if !p.resolved {
		select {
		case data = <-p.ch_resolved_value:
		case err = <-p.ch_error:
		}
	}
	return data, err
}

// Resolves a promise by updating it with a value.
func (p *Promise[T]) resolve(value T) {
	if p.resolved {
		return
	}
	p.resolved = true
	defer close(p.ch_resolved_value)
	defer close(p.ch_error)
	p.ch_resolved_value <- value
}

// Rejects a promise by updating it with an error.
func (p *Promise[T]) reject(err error) {
	if p.resolved {
		return
	}
	p.resolved = true
	defer close(p.ch_resolved_value)
	defer close(p.ch_error)
	p.ch_error <- err
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

type promise_resolver[R any] func(func(R), func(error))

// New creates a new Promise with the provided resolver function.
func New[R any](resolver promise_resolver[R]) *Promise[R] {
	promise := Promise[R]{
		ch_resolved_value: make(chan R),
		ch_error:          make(chan error),
		resolved:          false,
	}
	go func() {
		defer promise.panic_handler()
		resolver(promise.resolve, promise.reject)
	}()
	return &promise
}
