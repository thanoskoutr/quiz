package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Quiz struct {
	question, answer interface{}
}

func main() {
	// Parse command-line flag
	csvFilename := flag.String("csv", "problems.csv", "filename of the quiz csv")
	flag.Parse()

	// Open CSV
	file, err := os.Open(*csvFilename)
	if err != nil {
		log.Fatal(err)
	}
	// Close file
	defer file.Close()

	//  Read CSV
	problemsCsv := csv.NewReader(file)

	// Read CSV line by line
	quizSlice := make([]Quiz, 0)
	problemCount := 0

	for {
		record, err := problemsCsv.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// Save question/answer
		question := strings.TrimSpace(record[0])
		answer := strings.TrimSpace(record[1])
		quiz := Quiz{question: question, answer: answer}

		quizSlice = append(quizSlice, quiz)

		problemCount++
	}

	// Print questions, get answers and keep score
	score := 0
	var userAnswer string

	for i, quiz := range quizSlice {
		// Print question
		fmt.Printf("Problem #%v: %s = ", i+1, quiz.question)

		// Read answer
		fmt.Scanln(&userAnswer)
		userAnswer = strings.TrimSpace(userAnswer)

		// Check answer
		if userAnswer == quiz.answer {
			score++
		}
	}

	// Print score
	fmt.Printf("You scored %v out of %v\n", score, problemCount)
}
