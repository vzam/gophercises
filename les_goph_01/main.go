package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	fileName := flag.String("file", "problems.csv", "filename of a csv file which contains the questions and answers of the quiz")
	limit := flag.Int("limit", 30, "time limit after which the quiz will stop")
	flag.Parse()

	fmt.Printf("Timeout in %d seconds\n", *limit)

	fileReader, err := os.Open(*fileName)
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(fileReader)
	problemReader := NewReaderFromCsv(csvReader)

	problems, err := problemReader.ReadAll()
	if err != nil {
		panic(err)
	}

	done := make(chan int)
	correctAnswers := 0

	go func() {
		for i, problem := range problems {
			_, err := fmt.Printf("Problem #%d: %v = ", i, problem.Question)
			if err != nil {
				panic(err)
			}

			correct, err := scanAnswer(problem)
			if err != nil {
				panic(err)
			} else if correct {
				correctAnswers++
			}
		}
		done <- 1
	}()

	timeout := time.After(time.Duration(*limit) * time.Second)

	select {
	case <-done:
		fmt.Printf("\nYou scored %d out of %d", correctAnswers, len(problems))

		// this will leak the go routine but since we exit the program anyway, this fine. see https://stackoverflow.com/a/50798441/3946803
		os.Exit(0)
	case <-timeout:
		fmt.Printf("\nTimeout expired. You scored %d out of %d", correctAnswers, len(problems))
	}
}

func scanAnswer(problem *Problem) (correct bool, err error) {
	var answer string
	_, err = fmt.Scanln(&answer)
	if err != nil {
		return false, err
	}

	return normalize(answer) == normalize(problem.Answer), nil
}

func normalize(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}
