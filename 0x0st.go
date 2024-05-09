package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
)

func uploadFile(url string, filePath string, expire string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Prepare the multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	// Add other fields
	_ = writer.WriteField("expire", expire)
	_ = writer.WriteField("secret", "")

	err = writer.Close()
	if err != nil {
		return "", err
	}

	// Create the request
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

func main() {
	// Read from clipboard
	text, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println("Failed to read from clipboard:", err)
		return
	}

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "example.*.txt")
	if err != nil {
		fmt.Println("Failed to create a temporary file:", err)
		return
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file afterwards

	_, err = tmpFile.WriteString(text)
	if err != nil {
		fmt.Println("Failed to write to temporary file:", err)
		return
	}
	tmpFile.Close()

	// Upload the file
	url := "https://0x0.st"
	expire := "6" // 6 hours
	response, err := uploadFile(url, tmpFile.Name(), expire)
	if err != nil {
		fmt.Println("Failed to upload file:", err)
		return
	}

	// Output the response
	fmt.Println("Uploaded to:", response)

	// Copy the result URL back to the clipboard
	err = clipboard.WriteAll(response)
	if err != nil {
		fmt.Println("Failed to copy URL to clipboard:", err)
		return
	}
	fmt.Println("URL copied to clipboard.")
}
