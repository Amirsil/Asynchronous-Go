package main

import (
	"reflect"
	"sync"
)

func Async(function interface{}, args ...interface{}) Awaitable {
	channel := make(chan interface{}, 1)

	go func() {
		defer close(channel)
		reflectedArguments := make([]reflect.Value, len(args))

		for index, arg := range args {
			reflectedArguments[index] = ReflectType(arg)
		}

		channel <- ReflectFunction(function).
			Call(reflectedArguments[:])[0].
			Interface()
	}()

	return channel
}

func AwaitAll(awaitables ...Awaitable) []interface{} {
	results := make([]interface{}, len(awaitables))
	wg := new(sync.WaitGroup)

	for index, awaitable := range awaitables {
		wg.Add(1)
		// index and awaitable could change before the goroutine is over
		staticIndex, staticAwaitable := index, awaitable
		go func() {
			defer wg.Done()
			results[staticIndex] = staticAwaitable.Await()
		}()

	}

	wg.Wait()
	return results
}

func CallWhenDone(function interface{}, awaitable Awaitable) {
	go func() {
		result := ReflectType(awaitable.Await())
		ReflectFunction(function).Call([]reflect.Value{result})
	}()
}

func CallWhenAllDone(function interface{}, awaitables ...Awaitable) {
	go func() {
		results := ReflectType(AwaitAll(awaitables...))
		ReflectFunction(function).Call([]reflect.Value{results})
	}()
}
