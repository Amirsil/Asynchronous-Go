package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/wesovilabs/koazee"
)

func main() {
	fmt.Println("Starting Execution")
	rand.Seed(time.Now().UnixNano())
	asyncAwaitTest()
	compareTimeDifference(awaitAllTest, SynchronousTest, 20)
	whenDoneTest()
	whenAllDoneTest()
}

func asyncAwaitTest() {
	fmt.Println("\nTesting async/await functionality")
	fmt.Printf("%v\n", Async(getJokeFromAPI).Await())
	fmt.Printf("%v\n", Async(randWithDelay, 2).Await())
}

func awaitAllTest(numberOfExecutions int) {
	awaitables := koazee.StreamOf(
		make([]string, numberOfExecutions)).
		Map(
			func(string) Awaitable {
				return Async(getJokeFromAPI)
			}).
		Do().Out().
		Val().([]Awaitable)

	for index, result := range AwaitAll(awaitables...) {
		fmt.Printf("%v: %v\n", index+1, result.(string))
	}
}

func SynchronousTest(numberOfExecutions int) {
	results := koazee.StreamOf(make([]int, numberOfExecutions)).
		Map(
			func(int) string {
				return getJokeFromAPI()
			}).
		Do().Out().
		Val().([]string)

	for index, result := range results {
		fmt.Printf("%v: %v\n", index+1, result)
	}
}

func whenDoneTest() {
	fmt.Println("\nTesting whenDone callback functionality")
	done := false

	CallWhenDone(
		func(randNum int) {
			fmt.Printf("%v\n", randNum)
			done = true
		}, Async(randWithDelay, 3))

	for !done {
	}
}

func whenAllDoneTest() {
	fmt.Println("\nTesting whenAllDone callback functionality")
	done := false

	CallWhenAllDone(
		func(results []interface{}) {
			for _, result := range results {
				fmt.Printf("result: %v\n\n", result)
			}
			done = true
		},
		Async(randWithDelay, 2),
		Async(getJokeFromAPI),
		Async(getJokeFromAPI),
	)

	for !done {
	}
}

func compareTimeDifference(
	concurrentFunc interface{},
	nonConcurrentFunc interface{},
	numberOfExecutions int) {

	conurrectTime := MeasureTime(concurrentFunc, numberOfExecutions)
	nonConcurrentTime := MeasureTime(nonConcurrentFunc, numberOfExecutions)
	fmt.Printf("Concurrent took %v and Normal took %v", conurrectTime, nonConcurrentTime)
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
