package main

import (
	"reflect"
	"sync"
)

func async(function interface{}, args ...interface{}) <-chan interface{} {
	promise := make(chan interface{}, 1)
	callableFunction := reflect.ValueOf(function)

	var injectableArguments []reflect.Value

	for _, arg := range args {
		correctTypeArg := reflect.New(reflect.TypeOf(arg)).Elem()
		correctTypeArg.Set(reflect.ValueOf(arg))
		injectableArguments = append(injectableArguments, correctTypeArg)
	}

	go func() {
		defer close(promise)
		promise <- callableFunction.Call(injectableArguments[:])[0].Interface()
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

func whenDone(awaitable <-chan interface{}, function interface{}) {
	go func() {
		callableFunction := reflect.ValueOf(function)
		item := <-awaitable
		arg := reflect.New(reflect.TypeOf(item)).Elem()
		arg.Set(reflect.ValueOf(item))
		callableFunction.Call([]reflect.Value{arg})
	}()
}
