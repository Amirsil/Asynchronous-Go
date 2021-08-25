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

func main() {
	rand.Seed(time.Now().UnixNano())
	asyncAwaitTest()
	awaitAllTest()
	whenDoneTest()
	whenAllDoneTest()
}

func asyncAwaitTest() {
	fmt.Println("\nTesting async/await functionality")
	fmt.Printf("%v\n", await(async(getJokeFromAPI)))
	fmt.Printf("%v\n", await(async(randWithDelay, 2)))
}

func awaitAllTest() {
	fmt.Print("\nTesting Asyncronous ")
	randNum := async(randWithDelay, 3)
	joke := async(getJokeFromAPI)
	fmt.Print("CPU Efficient Processing\n")

	for index, result := range awaitAll(randNum, joke) {
		fmt.Printf("%v: %v\n", index+1, result)
	}
}

func whenDoneTest() {
	fmt.Println("\nTesting whenDone callback functionality")
	done := false

	whenDone(
		func(randNum int) {
			fmt.Printf("%v\n", randNum)
			done = true
		}, async(randWithDelay, 3))

	for !done {
	}
}

func whenAllDoneTest() {
	fmt.Println("\nTesting whenAllDone callback functionality")
	done := false

	whenAllDone(
		func(results []interface{}) {
			fmt.Printf("%v\n", results)
			done = true
		},
		async(getJokeFromAPI),
		async(getJokeFromAPI),
		async(getJokeFromAPI),
		async(getJokeFromAPI),
		async(getJokeFromAPI),
		async(getJokeFromAPI),
		async(getJokeFromAPI),
		async(getJokeFromAPI),
	)

	for !done {
	}
}

func randWithDelay(delay int) int {
	time.Sleep(time.Duration(delay) * time.Second)
	return rand.Int()
}

func getJokeFromAPI() string {
	httpResp, err := http.Get("https://v2.jokeapi.dev/joke/Any")
	if err != nil {
		log.Fatal(err)
	}

	jsonResponse, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response map[string]interface{}
	json.Unmarshal([]byte(string(jsonResponse)), &response)

	var joke string
	if _, ok := response["joke"]; ok {
		joke = fmt.Sprintf("\njoke: %v", response["joke"])
	} else {
		joke = fmt.Sprintf("\nquestion: %v\nanswer: %v", response["setup"], response["delivery"])
	}

	return joke
}
