package reviewbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const geminiURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key="

func CallGemini(prompt, apiKey string) (string, error) {
	type contentPart struct {
		Text string `json:"text"`
	}
	type content struct {
		Parts []contentPart `json:"parts"`
	}
	type requestBody struct {
		Contents []content `json:"contents"`
	}

	body := requestBody{
		Contents: []content{
			{
				Parts: []contentPart{{Text: prompt}},
			},
		},
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", geminiURL+apiKey, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Gemini API error: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal(raw, &result)
	if err != nil {
		return "", fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	candidates := result["candidates"].([]interface{})
	parts := candidates[0].(map[string]interface{})["content"].(map[string]interface{})["parts"].([]interface{})
	return parts[0].(map[string]interface{})["text"].(string), nil
}
