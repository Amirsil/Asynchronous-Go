package main

import (
	"fmt"
	"time"
)

func infinity() <-chan int {
	pipe := make(chan int, 1)

	go func() {
		defer close(pipe)
		num := 0
		for {
			pipe <- num
			num++
			time.Sleep(500 * time.Millisecond)
		}

	}()

	return pipe
}

func square(input <-chan int) chan int {
	pipe := make(chan int, cap(input))

	go func() {
		defer close(pipe)
		for item := range input {
			pipe <- item * item
		}
	}()

	return pipe
}

func main() {
	for i := range square(infinity()) {
		fmt.Printf("%v\n", i)
	}
}
