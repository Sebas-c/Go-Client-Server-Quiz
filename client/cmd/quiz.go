/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

type Question struct {
	Number   int
	Question string
	Answers  []string
}

// quizCmd represents the quiz command
var quizCmd = &cobra.Command{
	Use:   "quiz",
	Short: "Starts a quiz",
	Long:  `Starts a quiz, answer all the questions and see how well you fared compare to our users!`,
	Run: func(cmd *cobra.Command, args []string) {
		postAnswers(questionUser(queryQuestions()))
	},
}

// Returns the questions as array
func queryQuestions() []Question {
	var AllQuestions []Question
	// Sending a request and unMarshalling the results
	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		fmt.Println("Whoopsie daisy, an error occured. We'll be back shortly!")
		fmt.Println("Error: ")
		fmt.Println(err)
	} else {
		fmt.Println("The game is afoot!")
		bytes, errRead := ioutil.ReadAll(resp.Body)
		if errRead != nil {
			fmt.Println("An error occured:", errRead)
		}
		errUnmarshal := json.Unmarshal(bytes, &AllQuestions)
		if errUnmarshal != nil {
			fmt.Println("An error occured:", errUnmarshal)
		}
	}
	defer resp.Body.Close()
	return AllQuestions
}

func questionUser(AllQuestions []Question) []int {
	var allAnswers []int
	for i := 0; i < len(AllQuestions); i++ {
		// Printing question followed by possible answers
		fmt.Printf("--- Question %d/%d: ---\n", AllQuestions[i].Number, len(AllQuestions))
		fmt.Println(AllQuestions[i].Question)
		fmt.Printf("\n--- Pick a single answer from the following %d: ---\n", len(AllQuestions[i].Answers))
		for j := 0; j < len(AllQuestions[i].Answers); j++ {
			fmt.Println(j+1, "-", AllQuestions[i].Answers[j])
		}
		fmt.Printf("Your answer: ")
		// Quick input validation
		for {
			var tempAnswer string
			fmt.Scanln(&tempAnswer)
			if answer, err := strconv.Atoi(tempAnswer); err == nil {
				if answer < len(AllQuestions[i].Answers)+1 {
					// Storing answer
					allAnswers = append(allAnswers, answer)
					break
				}
			}
			fmt.Printf("Valid answers are contained between %d to %d (numerical values only)\n", 1, len(AllQuestions[i].Answers))
		}
	}

	// Prints all the answers for debugging
	// fmt.Println(allAnswers)
	return allAnswers
}

// Sends the answers to the server and displays the response to the user
func postAnswers(allAnswers []int) {
	var results []int

	// Marshalling answers to send
	JSONAnswers, JSONerr := json.Marshal(allAnswers)
	if JSONerr != nil {
		fmt.Print("Error while encoding the answers")
	}
	// Posting results and unmarshalling response
	response, errResults := http.Post("http://localhost:8080/answers", "application/json", bytes.NewBuffer(JSONAnswers))
	if errResults != nil {
		fmt.Println("An error occured:", errResults)
	}
	bytes, errRead := ioutil.ReadAll(response.Body)
	if errRead != nil {
		fmt.Println("An error occured:", errRead)
	}
	errUnmarshal := json.Unmarshal(bytes, &results)
	if errUnmarshal != nil {
		fmt.Println("An error occured:", errUnmarshal)
	}
	defer response.Body.Close()

	fmt.Printf("\n--- RESULTS ---\n")

	// Switch case based on the results obtained (Not necessary to go for a switch here,
	// just wanted to make it a bit more customized for the end user)
	switch results[0] {
	case 0:
		fmt.Printf("Though luck, you scored %v%%, that's better than %v%% users :( Better luck next time!", results[0], results[1])
	case 33:
		fmt.Printf("You scored %v%%, that's better than %v%% users. You can do better though!", results[0], results[1])
	case 66:
		fmt.Printf("Good job, you scored %v%%, that's better than %v%% users! Close to the top!", results[0], results[1])
	case 100:
		fmt.Printf("Awesome! You scored %v%%, that's better than %v%% users! Congrats!", results[0], results[1])
	}

}

func init() {
	rootCmd.AddCommand(quizCmd)
}
