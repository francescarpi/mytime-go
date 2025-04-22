package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mytime/settings"
	"mytime/utils"
	"net/http"

	"gorm.io/gorm"
)

type RedmineIssueProject struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type RedmineIssue struct {
	Id      int                 `json:"id"`
	Project RedmineIssueProject `json:"project"`
}

type RedmineIssueResponse struct {
	Issue RedmineIssue `json:"issue"`
}

type RedmineProjectActivity struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type RedmineProject struct {
	Id                  int                      `json:"id"`
	TimeEntryActivities []RedmineProjectActivity `json:"time_entry_activities"`
}

type RedmineProjectResponse struct {
	Project RedmineProject `json:"project"`
}

type RedmineSendTaskError struct {
	Errors []string `json:"errors"`
}

type RedmineConfig struct {
	DefaultActivity string `json:"default_activity"`
	Token           string `json:"token"`
	Url             string `json:"url"`
}

type Redmine struct {
	Config *RedmineConfig
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////
func requestPOST[T any](token, url string, rqBody io.Reader) (*T, error) {
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

// ////////////////////////////////////////////////////////////////////////////////////////////////////////
func requestGET[T any](token, url string) (*T, error) {
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

// ////////////////////////////////////////////////////////////////////////////////////////////////////////
func (r *Redmine) GetIssue(externalId string) (*RedmineIssue, error) {
	url := fmt.Sprintf("%s/issues/%s.json", r.Config.Url, externalId)

	response, err := requestGET[RedmineIssueResponse](r.Config.Token, url)
	if err != nil {
		log.Println("Error fetching issue:", err)
		return nil, err
	}
	return &response.Issue, nil
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////
func (r *Redmine) LoadActivities(externalId string) (*[]RedmineProjectActivity, error) {
	issue, err := r.GetIssue(externalId)
	if err != nil {
		log.Println("Error getting issue:", err)
		return nil, err
	}

	url := fmt.Sprintf("%s/projects/%d.json?include=time_entry_activities", r.Config.Url, issue.Project.Id)
	response, err := requestGET[RedmineProjectResponse](r.Config.Token, url)
	if err != nil {
		log.Println("Error fetching project:", err)
		return nil, err
	}
	return &response.Project.TimeEntryActivities, nil
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////
func (r *Redmine) SendTask(externalId, desc, date string, duration float64, activityId int) error {
	url := fmt.Sprintf("%s/time_entries.json", r.Config.Url)

	jsonBody := fmt.Sprintf(`{
		"time_entry": {
			"issue_id": "%s",
			"hours": "%s",
			"comments": "%s",
			"spent_on": "%s",
			"activity_id": %d
		}
	}`,
		externalId,
		utils.HumanizeDuration(duration),
		desc,
		date,
		activityId,
	)

	body := []byte(jsonBody)

	response, err := requestPOST[RedmineSendTaskError](r.Config.Token, url, bytes.NewBuffer(body))
	if response == nil && err == nil {
		// All good
		return nil
	}

	if response == nil && err != nil {
		// Generic error
		return err
	}

	log.Println("Task not sent successfully:", response)

	return fmt.Errorf("Error sending task: %v", response.Errors)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////
func NewRedmine(conn *gorm.DB) *Redmine {
	settings, err := settings.GetSettings(conn)

	if err != nil {
		log.Println("Error loading settings:", err)
		return nil
	}

	var config RedmineConfig
	err = json.Unmarshal([]byte(settings.IntegrationConfig), &config)
	if err != nil {
		log.Println("Error unmarshalling Redmine config:", err)
		return nil
	}

	return &Redmine{
		Config: &config,
	}
}
