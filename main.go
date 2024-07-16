package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// Question struct represents a question with its answer and explanation
type Question struct {
	Question     string
	Answer       string
	Explanation  string
}

// readQuestionsFromFile reads questions from  a local text file
func readQuestionsFromFile(filePath string) ([]Question, error) {
	// reaf file content
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil,  err
    }

	// split the file content into lines
	lines := strings.Split(string(file), "\n")
	questions := []Question{}

	// parse each line into a question struct
	for _, line := range lines {
		parts := strings.Split(line, "|") // assuming | as delimeter in text file
		if len(parts) == 3 {
			questions = append(questions, Question{Question: parts[0], Answer: parts[1], Explanation: parts[2]})
		}
	}
	return questions, nil 
}

// writeQuestionsToCSV writes questions to a CSV file
func writeQuestionsToCSV(questions []Question, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write the CSV header
	writer.Write([]string{"Question", "Answer", "Explanation"})

	// write each question to CSV file
	for _, question := range questions {
		writer.Write([]string{question.Question, question.Answer, question.Explanation})

	}
	return nil
}