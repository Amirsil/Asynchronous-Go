package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func async(f func() interface{}) <-chan interface{} {
	pipe := make(chan interface{}, 1)

	go func() {
		defer close(pipe)
		pipe <- f()
	}()

	return pipe
}

func await(awaitable <-chan interface{}) interface{} {
	return <-awaitable
}

func getNumber() interface{} {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(5 * time.Second)
	return rand.Int()
}

func getJoke() interface{} {
	httpResp, err := http.Get("https://v2.jokeapi.dev/joke/Any")
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	var result map[string]interface{}
	json.Unmarshal([]byte(string(bodyBytes)), &result)
	return fmt.Sprintf("%v\nanswer: %v", result["setup"], result["delivery"])
}

func main() {
	randNum := async(getNumber)
	joke := async(getJoke)

	fmt.Printf("%v\n", await(joke))
	fmt.Printf("random number: %v\n", await(randNum))
}
