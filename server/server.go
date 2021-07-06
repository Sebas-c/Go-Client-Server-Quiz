package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Struct for the questions. "Number" could be used as an ID for each question if the quiz became more evolved
type question struct {
	Number   int
	Question string
	Answers  []string
}

// Considering this is just an example, the correct answers and global results are hardcoded.
var correctAnswers = [3]int{1, 1, 2}

// Saved results as a map instead of an array/slice of all results. May not look like the most obvious choice
// but with a larger amount of users it would save computer cycles when checking results (O(4) vs O(n))
// Eventually this would depend on the DB's design anyway
var globalResults = map[int]int{
	0:   18,
	33:  33,
	66:  58,
	100: 26,
}

// Sends a set of questions
func homePage(w http.ResponseWriter, r *http.Request) {
	AllQuestions := []question{{1, "Who plays Mr. Robot?", []string{"Christian Slater", "Rami Malek", "Portia Doubleday", "Carly Chaikin"}},
		{2, "A boolean holds only a True or False value?", []string{"True", "False"}},
		{3, "Which word contains 6 letters?", []string{"Alphabet", "Goggle", "Pendelum"}}}
	JSONQuestions, JSONerr := json.Marshal(AllQuestions)
	if JSONerr != nil {
		fmt.Print("Error while encoding JSONQuestions")
	}

	fmt.Fprintf(w, string(JSONQuestions))
	// Printing request on main page if needed for debugging
	// fmt.Println("Questions sent to a client")
}

// Retrieves answers
func answers(w http.ResponseWriter, r *http.Request) {
	// Variables
	var answers []int
	correctResults := 0.0
	resultsBetterThan := 0
	totalUsers := 0

	// Retrieving results as JSON and unmarshalling
	bytes, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		fmt.Println("An error occured:", errRead)
	}
	errUnmarshal := json.Unmarshal(bytes, &answers)
	if errUnmarshal != nil {
		fmt.Println("An error occured:", errUnmarshal)
	}

	// Counting the correct answers
	for i := 0; i < len(answers); i++ {
		if correctAnswers[i] == answers[i] {
			correctResults++
		}
	}
	// Converting result as a percentage
	correctResults = correctResults * 100 / 3

	// Incrementing global results with user's results
	globalResults[int(correctResults)]++

	// Calculating how well the user did compared to others
	for score, numberOfUsers := range globalResults {
		if score < int(correctResults) {
			resultsBetterThan += numberOfUsers
		}
		totalUsers += numberOfUsers
	}
	// Converting to percentage
	resultsBetterThan = resultsBetterThan * 100 / totalUsers

	sendingResults(w, r, correctResults, resultsBetterThan)
	// Printing results if needed for debugging
	// fmt.Printf("User scored: %v. Those results were better than %v%%\n", correctResults, resultsBetterThan)
}

// Sends results as JSON back to the user
func sendingResults(w http.ResponseWriter, r *http.Request, userScore float64, placement int) {
	results := []int{int(userScore), placement}
	JSONQuestions, JSONerr := json.Marshal(results)
	if JSONerr != nil {
		fmt.Print("Error while encoding JSONQuestions")
	}

	fmt.Fprintf(w, string(JSONQuestions))
}

// Request handler
func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/answers", answers)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	fmt.Println("Server started")
	handleRequests()
}
