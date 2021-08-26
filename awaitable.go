package main

type Awaitable (<-chan interface{})

func (awaitable Awaitable) Await() interface{} {
	return <-awaitable
}
