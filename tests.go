package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Starting Execution")
	rand.Seed(time.Now().UnixNano())

	errorHandlingTest()
	time.Sleep(1 * time.Second)
	// compareTimeDifference(ConcurrentTest, SynchronousTest, 20)
}

func errorHandlingTest() {
	joke, _ := async[string](getJokeFromAPI).await()
	fmt.Printf("```%T```: %v\n", joke, joke)

	_, err := async[int](randWithDelay, 2).await()
	fmt.Printf("%v\n", err)

	_, err = async[int](randWithDelay, 7).await()

	fmt.Printf("%v\n", err)

	fmt.Println("\nTested async/await functionality")
}

func randWithDelay(delay int) (int, error) {
	if delay > 5 {
		return 0, errors.New("Delay is too big!")
	}

	time.Sleep(time.Duration(delay) * time.Second)
	return rand.Int(), nil
}

func getJokeFromAPI() (string, error) {
	httpResp, err := http.Get("https://v2.jokeapi.dev/joke/Any")
	if err != nil {
		return "", err
	}

	jsonResponse, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return "", err
	}

	var response map[string]interface{}
	json.Unmarshal([]byte(string(jsonResponse)), &response)
	var joke string

	if _, ok := response["joke"]; ok {
		joke = fmt.Sprintf("\njoke: %v", response["joke"])
	} else {
		joke = fmt.Sprintf("\nquestion: %v\nanswer: %v", response["setup"], response["delivery"])
	}

	return joke, nil
}
