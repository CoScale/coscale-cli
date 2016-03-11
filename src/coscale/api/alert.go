package api

import (
	"fmt"
)

type Alert struct {
	ID                int64
	Name              string
	Description       string
	Handle            string
	BackupSeconds     int
	BackupHandle      string
	EscalationSeconds int
	EscalationHandle  string
	Version           int64
	Source            string
}

func (e Alert) GetId() int64 {
	return e.ID
}

//GetAlertsBy will use a custom query to get a alert by unresolved/unacknowledged
func (api *Api) GetAlertsBy(query string) (string, error) {
	var result string
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/alerts/?%s=false", api.AppID, query), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

//AlertSolution will be used to acknowledge/ resolve a alert
func (api *Api) AlertSolution(alert *Alert, solutionType string) (string, error) {
	data := map[string][]string{
		"version": {fmt.Sprintf("%d", alert.Version)},
	}
	var result string
	if err := api.makeCall("PUT", fmt.Sprintf("/api/v1/app/%s/alerts/%d/%s/", api.AppID, alert.ID, solutionType), data, true, &result); err != nil {
		return "", err
	}
	return result, nil
}
