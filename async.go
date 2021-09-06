package main

import (
	"reflect"
)

func async[T any](function interface{}, args ...interface{}) awaitable[T] {
	valueChannel := make(chan T, 1)
	errorChannel := make(chan error, 1)

	go func() {
		defer close(valueChannel)
		defer close(errorChannel)
		reflectedArguments := make([]reflect.Value, len(args))

		for index, arg := range args {
			reflectedArguments[index] = ReflectType(arg)
		}

		result := ReflectFunction(function).Call(reflectedArguments[:])
		valueChannel <- result[0].Interface().(T)
		errorChannel <- GetFunctionError(result)
	}()

	return awaitable[T]{valueChannel, errorChannel}
}
