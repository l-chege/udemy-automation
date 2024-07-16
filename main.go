package main

import (
	"encoding/csv"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// question struct represents a question with its answer and explanation
type Question struct {
	Question     string
	Answer       string
	Explanation  string
}

//downloadFileFromBlob downloads a file from Azure Blob storage
func downloadFileFromBlob(accountName, accountKey, containerName, filePath string) error {
	// create blob url
	url := fmt.Sprintf("https.//%s.blob.core.windows.net/%s", accountName, containerName)

	// create a new shared key credential
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}

	// create a new blob client
	client, err := azblob.NewBlobClientWithSharedKeyCredential(url, cred, nil)
	if err != nil {
		return err
	}

	// download blob content
	ctx := context.Background()
	downloadResponse, err := client.Download(ctx, nil)
	if err != nil {
		return err
	}
	bodyStream := downloadResponse.Body(nil)

	// read blob content
	downloadData, err := ioutil.ReadAll(bodyStream)
	if err != nil {
		return err
	}

	// write content to local file
	return ioutil.WriteFile(downloadPath, downloadData, 0644)
}

// uploadFileToBlob uploads a file to Azure Blob Storage
func uploadFileToBlob(accountName, accountKey, containerName, filePath string) error {
	// create the Blob URL
	url := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	// create a new shared key credential
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}

	// create a new Container client
	containerClient, err := azblob.NewContainerClientWithSharedKeyCredential(url, cred, nil)
	if err != nil {
		return err
	}

	// open the local file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// upload the file to Blob Storage
	ctx := context.Background()
	_, err = containerClient.UploadFileToBlockBlob(ctx, file, azblob.UploadOption{})
	return err
}

// readQuestionsFromFile reads questions from a local text file
func readQuestionsFromFile(filePath string) ([]Question, error) {
	// read file content
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	
	// split file content into lines
	lines := strings.Split(string(file), "\n")
	questions := []Question{}

	//parse each line into a question struct
	for _, line := range lines {
		parts := strings.Split(line, "|") // assuming '|' as the delimiter in text file
		if len(parts) == 3 {
			questions = append(questions, Question{Question: parts[0], Answer: parts[1], Explanation: parts[2]})
		}
	}
	return questions, nil
}

// writeQuestionsToCSV writes questions to a CSV file
func writeQuestionsToCSV(questions []Question, outputPath string) error {
	// create output CSV file
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// create new csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write csv header
	writer.Write([]string{"Question", "Answer", "Explanation"})

	// Write each question to the CSV file
	for _, question := range questions {
		writer.Write([]string{question.Question, question.Answer, question.Explanation})
	}
	return nil
}

func main() {
	// Define Azure Blob Storage credentials and parameters
	accountName := "your_account_name"
	accountKey := "your_account_key"
	containerName := "your_container_name"
	inputBlobName := "questions.txt"
	outputBlobName := "questions.csv"
	downloadPath := "input/questions.txt"
	outputPath := "output/questions.csv"

	// Download the input file from Blob Storage
	err := downloadFileFromBlob(accountName, accountKey, containerName, inputBlobName, downloadPath)
	if err != nil {
		log.Fatalf("Error downloading file: %v", err)
	}

	// Read questions from the downloaded file
	questions, err := readQuestionsFromFile(downloadPath)
	if err != nil {
		log.Fatalf("Error reading questions: %v", err)
	}

	// Write the questions to a CSV file
	err = writeQuestionsToCSV(questions, outputPath)
	if err != nil {
		log.Fatalf("Error writing to CSV: %v", err)
	}

	// upload the output csv file to blob storage
	err = UploadFileToBlob(accountName, accountKey, containerName, outputPath)
	if err != nil {
	    log.Fatal("Error uploading: %v", err)
	}

	fmt.Println("File processed and uploaded successfully!")
}

