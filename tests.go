package main

import (
	"os"
	"testing"
)

func TestDownloadFileFromBlob(t *testing.T) {
	accountName := "your_account_name"
	accountKey := "your_account_key"
	containerName := "your_container_name"
	inputBlobName := "questions.txt"
	downloadPath := "input/questions.txt"

	err := downloadFileFromBlob(accountName, accountKey, containerName, inputBlobName, downloadPath)
    if err != nil {
        t.Fatalf("Failed to download file: %v", err)
    }

	if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
		t.Fatalf("Downloaded file does not exist")
	}
}

func TestReadQuestionsFromFile(t *testing.T) {
	filePath := "input/questions.txt"
	questions, err := readQuestionsFromFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read questions: %v", err)
	}

	if len(questions) == 0 {
		t.Fatalf("No questions read from file")
	}
}

func TestWriteQuestionsToCSV(t *testing.T) {
    questions := []Question {
        {"What is Go?", "A programming language created by Google.", "Go is known for its simplicity and performance."},
        {"What is a goroutine?", "A lightweight thread managed by Go runtime.", "Goroutines are cheaper than threads."},
    }
	outputPath := "output/questions.csv"
	err := writeQuestionsToCSV(questions, outputPath)
	if err != nil {
		t.Fatalf("Failed to write to CSV: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("CSV file does not exist")
	}
}
