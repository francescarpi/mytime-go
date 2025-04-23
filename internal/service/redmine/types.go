package redmine

import (
	"encoding/json"
	"fmt"
	"strconv"
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
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Default bool
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

type IntFromString int

func (i *IntFromString) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case float64:
		*i = IntFromString(int(v))
		return nil
	case string:
		n, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		*i = IntFromString(n)
		return nil
	default:
		return fmt.Errorf("invalid type for IntFromString: %T", v)
	}
}
