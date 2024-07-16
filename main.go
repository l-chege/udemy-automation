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

// Question struct represents a question with its answer and explanation
type Question struct {
	Question     string
	Answer       string
	Explanation  string
}

//downloadFileFromBlob downloads a file from Azure Blob storage
func downloadFileFromBlob(accountName, accountKey, containerName, filePath string) error {
	// create blob url
	url := fmt.Sprintf("https.//%s.blob.core.windows.net/%s", accountName, containerName)

	// create a new shard key credential
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}

	// create a new blob client
	client, err := azblob.NewBlobClientWithSharedKeyCredential(url, crd, nil)
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
	// Construct the Blob URL
	url := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	// Create a new shared key credential
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}

	// Create a new Container client
	containerClient, err := azblob.NewContainerClientWithSharedKeyCredential(url, cred, nil)
	if err != nil {
		return err
	}

	// Open the local file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Upload the file to Blob Storage
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
}





// readQuestionsFromFile reads questions from  a local text file
func readQuestionsFromFile(filePath string) ([]Question, error) {
	// read file content
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

// main function to orchestrate reading from input file and writing to output csv
func main() {
	inputFilePath := "input/questions.txt"
	outputFilePath := "output/questions.csv"

	questions, err := readQuestionsFromFile(inputFilePath)
	if err != nil {
		fmt.Println("Error reading questions:", err)
		return
	}

	err = writeQuestionsToCSV(questions, outputFilePath)
	if err != nil {
		fmt.Println("Error writing to CSV:", err)
	}
}