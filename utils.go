package main

import (
	"reflect"
	"time"
)

func ReflectFunction(function interface{}) reflect.Value {
	reflectedFunction := reflect.ValueOf(function)

	if reflectedFunction.Kind() != reflect.Func {
		panic("Expected function parameter is not a function")
	}

	return reflectedFunction
}

func ReflectType(value interface{}) reflect.Value {
	reflectedValue := reflect.New(reflect.TypeOf(value)).Elem()
	reflectedValue.Set(reflect.ValueOf(value))

	return reflectedValue
}

func MeasureTime(function interface{}, args ...interface{}) time.Duration {
	reflectedArguments := make([]reflect.Value, len(args))

	for index, arg := range args {
		reflectedArguments[index] = ReflectType(arg)
	}

	start := time.Now()
	ReflectFunction(function).Call(reflectedArguments[:])
	
	return time.Since(start)
}

func GetFunctionError(returnValue []reflect.Value) error {
	err := error(nil)

	if len(returnValue) > 1 {
		returnedError := returnValue[1].Interface()

		if returnedError != nil {
			err = returnedError.(error)
		}
	}

	return err
}
