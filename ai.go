package reviewbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	HuggingFaceAPIURL = "https://api-inference.huggingface.co/models/"
)

type RequestBody struct {
	Inputs string `json:"inputs"`
}

type ResponseBody struct {
	GeneratedText string `json:"generated_text"`
}

var httpClient = &http.Client{}

// CallHuggingFace interacts with Hugging Face's inference API to generate text from a model.
func GetReviews(modelName, apiToken, prompt string) (string, error) {
	// Construct the request body
	requestBody := RequestBody{Inputs: prompt}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create the POST request with headers for authorization
	req, err := http.NewRequest("POST", HuggingFaceAPIURL+modelName, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set the necessary headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("hugging Face API error: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("received non-OK response from Hugging Face API: %d %s", resp.StatusCode, string(body))
	}

	// Read and parse the response body
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response JSON into a ResponseBody struct
	var response []ResponseBody
	err = json.Unmarshal(raw, &response)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract and return the generated text
	if len(response) > 0 {
		return response[0].GeneratedText, nil
	}

	return "", fmt.Errorf("unexpected response structure: no generated text found")
}

// LoadHuggingFaceToken reads the Hugging Face API token from an environment variable
func LoadHuggingFaceToken() string {
	apiToken := os.Getenv("HUGGING_FACE_API_TOKEN")
	if apiToken == "" {
		log.Fatal("Hugging Face API token is not set in environment variables")
	}
	return apiToken
}
