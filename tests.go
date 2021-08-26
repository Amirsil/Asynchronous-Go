package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	Async(getJokeFromAPI).
		Then(func(item string) (int, error) {
			fmt.Printf("item: %v\n", item)
			return 0, errors.New("My Bad")
		}).
		Catch(func(err error) { log.Fatal(err) }).
		Await()

	// fmt.Println("Starting Execution")
	// rand.Seed(time.Now().UnixNano())
	// asyncAwaitTest()
	// compareTimeDifference(awaitAllTest, SynchronousTest, 20)
	// whenDoneTest()
	// whenAllDoneTest()
}

func asyncAwaitTest() {
	fmt.Println("\nTesting async/await functionality")
	joke, _ := Async(getJokeFromAPI).Await()
	fmt.Printf("%v\n", joke)

	_, err := Async(randWithDelay, 2).Await()
	fmt.Printf("%v\n", err)

	_, err = Async(randWithDelay, 7).Await()

	fmt.Printf("%v\n", err)
}

func randWithDelay(delay int) (int, error) {
	if delay > 5 {
		return 0, errors.New("Delay is too big!")
	}

	time.Sleep(time.Duration(delay) * time.Second)
	return rand.Int(), nil
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
