package main

import (
	"reflect"
)

type Awaitable struct {
	value <-chan interface{}
	err   <-chan error
}

func (awaitable Awaitable) Await() (interface{}, error) {
	return <-awaitable.value, <-awaitable.err
}

func (awaitable Awaitable) Then(callback interface{}) Awaitable {
	valueChannel := make(chan interface{}, 1)
	errorChannel := make(chan error, 1)

	go func() {
		defer close(valueChannel)
		defer close(errorChannel)
		value, _ := awaitable.Await()

		returnValue := ReflectFunction(callback).
			Call([]reflect.Value{ReflectType(value)})

		if len(returnValue) == 0 {
			valueChannel <- nil
		} else {
			valueChannel <- returnValue[0].Interface()
		}

		errorChannel <- GetFunctionError(returnValue)
	}()

	return Awaitable{valueChannel, errorChannel}
}

func (awaitable Awaitable) Catch(handleErr interface{}) Awaitable {
	valueChannel := make(chan interface{}, 1)
	errorChannel := make(chan error, 1)

	go func() {
		defer close(valueChannel)
		defer close(errorChannel)
		value, err := awaitable.Await()

		if err != nil {
			ReflectFunction(handleErr).
				Call([]reflect.Value{
					ReflectType(err)})
		}
		valueChannel <- value
		errorChannel <- nil
	}()

	return Awaitable{valueChannel, errorChannel}
}
