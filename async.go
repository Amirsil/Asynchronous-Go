package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func async(f func(...interface{}) interface{}, args ...interface{}) <-chan interface{} {
	pipe := make(chan interface{}, 1)

	go func() {
		defer close(pipe)
		pipe <- f(args...)
	}()
	return pipe
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

// randWithDelay(delay int) (int)
func randWithDelay(args ...interface{}) interface{} {
	delay := time.Duration(args[0].(int))

	rand.Seed(time.Now().UnixNano())
	time.Sleep(delay * time.Second)
	return rand.Int()
}

// getJokeFromAPI() (string)
func getJokeFromAPI(args ...interface{}) interface{} {
	httpResp, err := http.Get("https://v2.jokeapi.dev/joke/Any")
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	var response map[string]interface{}
	json.Unmarshal([]byte(string(bodyBytes)), &response)

	var joke string
	if _, ok := response["joke"]; ok {
		joke = fmt.Sprintf("\njoke: %v", response["joke"])
	} else {
		joke = fmt.Sprintf("\nquestion: %v\nanswer: %v", response["setup"], response["delivery"])
	}

	return joke
}

func main() {
	randNum := async(randWithDelay, 3)
	joke := async(getJokeFromAPI)

	results := awaitAll(randNum, joke)
	for index, result := range results {
		fmt.Printf("%v: %v\n", index+1, result)
	}

	time.Sleep(1)
	fmt.Printf("%v\n", await(async(getJokeFromAPI)))
	fmt.Printf("%v\n", await(async(randWithDelay, 2)))
}
