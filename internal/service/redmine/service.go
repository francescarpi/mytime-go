package redmine

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/francescarpi/mytime/internal/service"
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
