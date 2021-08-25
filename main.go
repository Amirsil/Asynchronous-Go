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

func randWithDelay(delay int) int {
	rand.Seed(time.Now().UnixNano())
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
