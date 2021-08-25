package main

import (
	"reflect"
	"sync"
)

func async(f interface{}, args ...interface{}) <-chan interface{} {
	promise := make(chan interface{}, 1)
	callableFunction := reflect.ValueOf(f)

	var injectableArguments []reflect.Value

	for _, arg := range args {
		correctTypeArg := reflect.New(reflect.TypeOf(arg)).Elem()
		correctTypeArg.Set(reflect.ValueOf(arg))
		injectableArguments = append(injectableArguments, correctTypeArg)
	}

	go func() {
		defer close(promise)
		promise <- callableFunction.Call(injectableArguments[:])[0]
	}()

	return promise
}

func await(awaitable <-chan interface{}) interface{} {
	return <-awaitable
}

func awaitAll(awaitables ...<-chan interface{}) []interface{} {
	results := make([]interface{}, len(awaitables))
	wg := new(sync.WaitGroup)

	for index, awaitable := range awaitables {
		wg.Add(1)
		// index and awaitable could change before the goroutine is over
		staticIndex, staticAwaitable := index, awaitable
		go func() {
			defer wg.Done()
			results[staticIndex] = <-staticAwaitable
		}()

	}
	wg.Wait()

	return results
}
