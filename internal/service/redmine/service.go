package redmine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/francescarpi/mytime/internal/service"
	"github.com/francescarpi/mytime/internal/util"
)

type Redmine struct {
	DefaultActivity IntFromString `json:"default_activity"`
	Token           string        `json:"token"`
	Url             string        `json:"url"`
}

func NewRedmine(service *service.Service) *Redmine {
	settings, err := service.Repo.GetSettings()
	if err != nil {
		panic(err)
	}

	var redmine Redmine
	err = json.Unmarshal([]byte(settings.IntegrationConfig), &redmine)
	if err != nil {
		panic(err)
	}

	return &redmine
}

func (r *Redmine) GetIssue(externalId string) (*RedmineIssue, error) {
	url := fmt.Sprintf("%s/issues/%s.json", r.Url, externalId)

	response, err := RequestGET[RedmineIssueResponse](r.Token, url)
	if err != nil {
		log.Println("Error fetching issue:", err)
		return nil, err
	}
	return &response.Issue, nil
}

func (r *Redmine) LoadActivities(externalId string) (*[]RedmineProjectActivity, *RedmineProjectActivity, error) {
	issue, err := r.GetIssue(externalId)
	if err != nil {
		log.Println("Error getting issue:", err)
		return nil, nil, err
	}

	url := fmt.Sprintf("%s/projects/%d.json?include=time_entry_activities", r.Url, issue.Project.Id)
	response, err := RequestGET[RedmineProjectResponse](r.Token, url)
	if err != nil {
		log.Println("Error fetching project:", err)
		return nil, nil, err
	}

	var defaultActivity RedmineProjectActivity
	activities := response.Project.TimeEntryActivities
	for _, activity := range activities {
		if activity.Id == int(r.DefaultActivity) {
			defaultActivity = activity
		}
	}

	return &activities, &defaultActivity, nil
}

func (r *Redmine) SendTask(externalId, desc, date string, duration int, activityId int) error {
	url := fmt.Sprintf("%s/time_entries.json", r.Url)

	type TimeEntry struct {
		IssueId    string `json:"issue_id"`
		Hours      string `json:"hours"`
		Comments   string `json:"comments"`
		SpentOn    string `json:"spent_on"`
		ActivityId int    `json:"activity_id"`
	}

	timeEntry := struct {
		TimeEntry TimeEntry `json:"time_entry"`
	}{
		TimeEntry: TimeEntry{
			IssueId:    externalId,
			Hours:      util.HumanizeDuration(duration),
			Comments:   desc,
			SpentOn:    date,
			ActivityId: activityId,
		},
	}

	body, err := json.Marshal(timeEntry)
	if err != nil {
		log.Println("Error marshalling time entry:", err)
		return err
	}

	response, err := RequestPOST[RedmineSendTaskError](r.Token, url, bytes.NewBuffer(body))
	if response == nil && err == nil {
		// All good
		return nil
	}

	if response == nil && err != nil {
		// Generic error
		log.Println("Error sending task:", err, bytes.NewBuffer(body))
		return err
	}

	log.Println("Task not sent successfully:", response.Errors)

	return fmt.Errorf("Error sending task: %v", response.Errors)
}
