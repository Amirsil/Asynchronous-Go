package main

import (
	"reflect"
)

func Async(function interface{}, args ...interface{}) Awaitable {
	valueChannel := make(chan interface{}, 1)
	errorChannel := make(chan error, 1)

	go func() {
		defer close(valueChannel)
		defer close(errorChannel)
		reflectedArguments := make([]reflect.Value, len(args))

		for index, arg := range args {
			reflectedArguments[index] = ReflectType(arg)
		}

		result := ReflectFunction(function).Call(reflectedArguments[:])

		valueChannel <- result[0].Interface()
		errorChannel <- GetFunctionError(result)
	}()

	return Awaitable{valueChannel, errorChannel}
}
