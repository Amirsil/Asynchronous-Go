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

	"github.com/wesovilabs/koazee"
)

func main() {
	fmt.Println("Starting Execution")
	rand.Seed(time.Now().UnixNano())

	errorHandlingTest()
	time.Sleep(1 * time.Second)
	promiseFunctionalityTest()
	time.Sleep(3 * time.Second)
	compareTimeDifference(ConcurrentTest, SynchronousTest, 20)
}

func errorHandlingTest() {
	joke, _ := Async(getJokeFromAPI).Await()
	fmt.Printf("%v\n", joke)

	_, err := Async(randWithDelay, 2).Await()
	fmt.Printf("%v\n", err)

	_, err = Async(randWithDelay, 7).Await()

	fmt.Printf("%v\n", err)

	fmt.Println("\nTested async/await functionality")
}

func promiseFunctionalityTest() {
	Async(randWithDelay, 4).
		Then(func(item int) (int, error) {
			fmt.Printf("item: %v\n", item)
			return 0, errors.New("My Bad")
		}).
		Catch(func(err error) { fmt.Printf("err: %v\n", err) }).
		Then(func(item int) { fmt.Printf("item: %v\n", item) }).
		Await()

	fmt.Println("\nTested Promise like functionality as in js")
}

func ConcurrentTest(numberOfExecutions int) {
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

func compareTimeDifference(
	concurrentFunc interface{},
	nonConcurrentFunc interface{},
	numberOfExecutions int) {

	conurrectTime := MeasureTime(concurrentFunc, numberOfExecutions)
	nonConcurrentTime := MeasureTime(nonConcurrentFunc, numberOfExecutions)
	fmt.Printf("Concurrent took %v and Normal took %v after %v runs",
		conurrectTime,
		nonConcurrentTime,
		numberOfExecutions)
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
