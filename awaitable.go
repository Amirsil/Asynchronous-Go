package main

type awaitable[T any] struct {
	value <-chan T
	err   <-chan error
}

func (awaitable awaitable[T]) await() (T, error) {
	return <-awaitable.value, <-awaitable.err
}

