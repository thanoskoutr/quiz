package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Quiz struct {
	question, answer interface{}
}

func main() {
	// Parse command-line flag
	csvFilename := flag.String("csv", "problems.csv", "filename of the quiz csv")
	timeLimit := flag.Int64("limit", 30, "time limit for each question of the quiz")
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

	// Create a variable for user's answer, a total score and a channel for receiving user's answer
	c := make(chan int)
	score := 0
	var userAnswer string

	// Print questions, get answers and keep score
	for i, quiz := range quizSlice {
		// Print question
		fmt.Printf("Problem #%v: %s = ", i+1, quiz.question)

		// Run getAnswer as a goroutine
		go getAnswer(userAnswer, quiz, score, c)

		// Wait for an answer, if timeLimit is exceeded stop the quiz
		select {
		case score = <-c:
			continue
		case <-time.After(time.Duration(*timeLimit) * time.Second):
			fmt.Println()
			fmt.Printf("You scored %v out of %v\n", score, problemCount)
			log.Fatal("Timed Out")
		}
	}

	// Print score
	fmt.Printf("You scored %v out of %v\n", score, problemCount)
}

func getAnswer(userAnswer string, quiz Quiz, score int, c chan int) {
	// Read answer
	fmt.Scanln(&userAnswer)
	userAnswer = strings.TrimSpace(userAnswer)

	// Check answer
	if userAnswer == quiz.answer {
		score++
	}
	c <- score
}
