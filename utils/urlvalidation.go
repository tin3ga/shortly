package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Base URL for MetaDefender API.
const metaDefenderAPI = "https://api.metadefender.com/v4/url/"

func URLValidation(link string, apiKey string) (string, error) {
	if apiKey == "" {
		// log.Print("API key missing")
		return "Api not found", fmt.Errorf("API key is required")
	}

	if link == "" {
		return "Empty Link", fmt.Errorf("Link URL is required")
	}

	encodedURL := url.QueryEscape(link)

	apiUrl := metaDefenderAPI + encodedURL

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		// log.Printf("Failed to create request: %v", err)
		return "Failed to create request", err
	}

	req.Header.Add("apiKey", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// log.Printf("Failed to perform request: %v", err)
		return "Failed to perform request", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		// log.Printf("Failed to read response body: %v", err)
		return "Failed to read response body", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		// log.Printf("API call failed with status: %d, response: %s", res.StatusCode, string(body))
		return fmt.Sprintf("API call failed with status: %d", res.StatusCode), fmt.Errorf("API error")
	}

	// Parse the JSON response body into a map.
	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		// log.Printf("Failed to parse response body: %v", err)
		return "Failed to parse response body", err
	}

	// responseMap := map[string]interface{}{
	// 	"address": "https://google.com",
	// 	"lookup_results": map[string]interface{}{
	// 		"detected_by": 0,
	// 		"sources": []map[string]interface{}{
	// 			{
	// 				"assessment":  "trustworthy",
	// 				"category":    "Search Engines",
	// 				"detect_time": "",
	// 				"provider":    "webroot.com",
	// 				"status":      0,
	// 				"update_time": "2025-01-09T08:59:38.413Z",
	// 			},
	// 		},
	// 		"start_time": "2025-01-09T09:47:56.759Z",
	// 	},
	// }
	lookupResults := responseMap["lookup_results"].(map[string]interface{})

	detectedBy := lookupResults["detected_by"].(float64)

	if detectedBy != 0.0 {
		return "malicious", nil

	} else {
		return "safe", nil

	}

}
