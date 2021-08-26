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
	callableFunction := ReflectFunction(function)
	reflectedArguments := make([]reflect.Value, len(args))

	for index, arg := range args {
		reflectedArguments[index] = ReflectType(arg)
	}

	start := time.Now()
	callableFunction.Call(reflectedArguments[:])
	return time.Since(start)
}