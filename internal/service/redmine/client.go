package redmine

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func RequestGET[T any](token, url string) (*T, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Redmine-API-Key", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error fetching issue:", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil, err
	}

	var result T
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		return nil, err
	}

	return &result, nil
}

func RequestPOST[T any](token, url string, rqBody io.Reader) (*T, error) {
	req, err := http.NewRequest("POST", url, rqBody)
	req.Header.Set("X-Redmine-API-Key", token)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error fetching issue:", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil, err
	}

	switch resp.StatusCode {
	case 201:
		return nil, nil
	case 401:
		return nil, fmt.Errorf("Unauthorized")
	}

	var result T
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		return nil, err
	}

	return &result, nil
}
