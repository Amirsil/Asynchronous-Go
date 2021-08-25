package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"sync"
	"time"
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

func randWithDelay(delay int) int {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(delay) * time.Second)
	return rand.Int()
}

func getJokeFromAPI() (string, error) {
	httpResp, err := http.Get("https://v2.jokeapi.dev/joke/Any")
	if err != nil {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", err
	}

	var response map[string]interface{}
	json.Unmarshal([]byte(string(bodyBytes)), &response)

	var joke string
	if _, ok := response["joke"]; ok {
		joke = fmt.Sprintf("\njoke: %v", response["joke"])
	} else {
		joke = fmt.Sprintf("\nquestion: %v\nanswer: %v", response["setup"], response["delivery"])
	}

	return joke, nil
}

func main() {
	fmt.Print("Starting Asyncronous ")
	randNum := async(randWithDelay, 3)
	joke := async(getJokeFromAPI)
	fmt.Print("CPU Efficient Work\n")

	results := awaitAll(randNum, joke)

	for index, result := range results {
		fmt.Printf("%v: %v\n", index+1, result)
	}

	time.Sleep(1)
	fmt.Printf("%v\n", await(async(getJokeFromAPI)))
	fmt.Printf("%v\n", await(async(randWithDelay, 2)))
}
