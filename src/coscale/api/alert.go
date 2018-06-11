package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Alert defines an alert for a trigger.
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

// GetId returns the id of the Alert.
func (e Alert) GetId() int64 {
	return e.ID
}

// AlertType defines a type for which Alerts could be inserted.
type AlertType struct {
	ID                int64
	Name              string
	Description       string
	Handle            string
	BackupHandle      string // Optional
	EscalationHandle  string // Optional
	BackupSeconds     int64  // Optional
	EscalationSeconds int64  // Optional
	Source            string
	Version           int64
}

// GetId returns the id of the AlertType.
func (a AlertType) GetId() int64 {
	return a.ID
}

// AlertTrigger defines what Triggers an Alert of AlertType
type AlertTrigger struct {
	ID             int64
	Name           string // Unique
	Description    string
	AutoResolve    int64  `json:"autoresolveSeconds"` // Optional
	DimensionSpecs string // The dimension specs for the selected metric.
	Metric         int64
	Config         string
	OnApp          bool
	GroupID        int64 // Optional
	ServerID       int64 // Optional
	Source         string
	Version        int64
}

// GetId returns the id of the AlertTrigger.
func (a AlertTrigger) GetId() int64 {
	return a.ID
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

// CreateType is used to add a new Alert type.
func (api *Api) CreateType(name, description, handle, backupHandle, escalationHandle string, backupSeconds, escalationSeconds int64) (string, error) {

	data := map[string][]string{
		"name":        {name},
		"description": {description},
		"handle":      {handle},
		"source":      {GetSource()},
	}

	// Set the optional values if they have value.
	if backupSeconds != -1 {
		data["backupSeconds"] = []string{fmt.Sprintf("%d", backupSeconds)}
	}
	if backupHandle != DEFAULT_STRING_VALUE {
		data["backupHandle"] = []string{backupHandle}
	}
	if escalationSeconds != -1 {
		data["escalationSeconds"] = []string{fmt.Sprintf("%d", escalationSeconds)}
	}
	if escalationHandle != DEFAULT_STRING_VALUE {
		data["escalationHandle"] = []string{escalationHandle}
	}

	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/alerttypes/", api.AppID), data, true, &result); err != nil {
		if duplicate, id := IsDuplicate(err); duplicate {
			return api.GetObject("alerttype", id)
		}
		return "", err
	}
	return result, nil
}

// UpdateType is used to update an existing Alert type.
func (api *Api) UpdateType(alertType *AlertType) (string, error) {

	data := map[string][]string{
		"name":        {alertType.Name},
		"description": {alertType.Description},
		"handle":      {alertType.Handle},
		"source":      {alertType.Source},
		"version":     {fmt.Sprintf("%d", alertType.Version)},
	}

	// Set the optional values if they have value.
	if alertType.BackupSeconds != 0 {
		data["backupSeconds"] = []string{fmt.Sprintf("%d", alertType.BackupSeconds)}
	}
	if alertType.BackupHandle != "" {
		data["backupHandle"] = []string{alertType.BackupHandle}
	}
	if alertType.EscalationSeconds != 0 {
		data["escalationSeconds"] = []string{fmt.Sprintf("%d", alertType.EscalationSeconds)}
	}
	if alertType.EscalationHandle != "" {
		data["escalationHandle"] = []string{alertType.EscalationHandle}
	}

	var result string
	if err := api.makeCall("PUT", fmt.Sprintf("/api/v1/app/%s/alerttypes/%d/", api.AppID, alertType.GetId()), data, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

// GetTriggers will return all triggers for an alert type
func (api *Api) GetTriggers(alertTypeID int64) (string, error) {
	var result string
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/alerttypes/%d/triggers/", api.AppID, alertTypeID), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

// CreateTrigger is used to add a new Trigger for alerts.
func (api *Api) CreateTrigger(name, description, config, dimensionSpecs string, alertTypeID, autoResolve, metricID, serverID, serverGroupID int64, onApp bool) (string, error) {

	data := map[string][]string{
		"name":        {name},
		"description": {description},
		"metric":      {fmt.Sprintf("%d", metricID)},
		"config":      {config},
		"onApp":       {fmt.Sprintf("%t", onApp)},
		"source":      {GetSource()},
	}

	parsedDimensionSpecs, err := ParseDimensionSpecs(dimensionSpecs)
	if err != nil {
		return "", err
	}

	data["dimensionSpecs"] = []string{parsedDimensionSpecs}

	// Set the option values if they have value.
	if serverID != -1 {
		data["server"] = []string{fmt.Sprintf("%d", serverID)}
	} else if serverGroupID != -1 {
		data["group"] = []string{fmt.Sprintf("%d", serverGroupID)}
	}
	if autoResolve != -1 {
		data["autoresolveSeconds"] = []string{fmt.Sprintf("%d", autoResolve)}
	}

	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/alerttypes/%d/triggers/", api.AppID, alertTypeID), data, true, &result); err != nil {
		if duplicate, id := IsDuplicate(err); duplicate {
			return api.GetObjectFromGroup("alerttype", "trigger", alertTypeID, id)
		}
		return "", err
	}
	return result, nil
}

// UpdateTrigger is used to update a existing Trigger for alerts.
func (api *Api) UpdateTrigger(typeID int64, trigger *AlertTrigger) (string, error) {

	data := map[string][]string{
		"name":           {trigger.Name},
		"description":    {trigger.Description},
		"dimensionSpecs": {trigger.DimensionSpecs},
		"config":         {trigger.Config},
		"onApp":          {fmt.Sprintf("%t", trigger.OnApp)},
		"source":         {trigger.Source},
		"version":        {fmt.Sprintf("%d", trigger.Version)},
	}

	// Set the option values if they have value.
	if trigger.Metric != 0 {
		data["metric"] = []string{fmt.Sprintf("%d", trigger.Metric)}
	}
	if trigger.ServerID != 0 {
		data["server"] = []string{fmt.Sprintf("%d", trigger.ServerID)}
	} else if trigger.GroupID != 0 {
		data["group"] = []string{fmt.Sprintf("%d", trigger.GroupID)}
	}
	if trigger.AutoResolve != 0 {
		data["autoresolveSeconds"] = []string{fmt.Sprintf("%d", trigger.AutoResolve)}
	}

	var result string
	if err := api.makeCall("PUT", fmt.Sprintf("/api/v1/app/%s/alerttypes/%d/triggers/%d/", api.AppID, typeID, trigger.ID), data, true, &result); err != nil {
		return "", err
	}
	return api.GetObjectFromGroup("alerttype", "trigger", typeID, trigger.ID)
}

// ParseHandle is used to parse the handle provided by user and serialize into json format.
func ParseHandle(handle string) (string, error) {
	var result []map[string]string

	// Parse the handle parameter.
	contacts := strings.Split(handle, " ")
	for _, contact := range contacts {
		// e.g. SLACK:slack/webhook/here
		if i := strings.Index(contact, ":"); i != -1 {
			contactType := contact[:i]
			contact = contact[i+1:]

			var contactRes map[string]string
			switch contactType {
			case "EMAILUSER":
				contactRes = map[string]string{
					"type": contactType,
					"id":   contact,
				}
			case "EMAIL":
				contactRes = map[string]string{
					"type":    contactType,
					"address": contact,
				}
			case "SLACK":
				contactRes = map[string]string{
					"type":    contactType,
					"webhook": contact,
				}
			}

			if contactRes != nil {
				result = append(result, contactRes)
			}
		}
	}
	if result == nil {
		return "", fmt.Errorf("Could not parse the alert handle")
	}

	jsonHandle, err := json.Marshal(result)
	return string(jsonHandle), err
}

// ParseDimensionSpecs is used to parse the dimensions specs for a metric.
func ParseDimensionSpecs(format string) (string, error) {
	var js interface{}
	// Check if the specs are already provided in JSON format.
	if err := json.Unmarshal([]byte(format), &js); err == nil {
		return format, nil
	}
	var buffer bytes.Buffer

	formatParts := strings.Split(format, ";")

	buffer.WriteString("[")
	for i, formatPart := range formatParts {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString("[")

		elements := strings.Split(formatPart, ":")

		var dimPart, dimValsPart, aggregatorPart string
		if len(elements) == 2 {
			dimPart = elements[0]
			dimValsPart = elements[1]
		} else if len(elements) == 3 {
			dimPart = elements[0]
			aggregatorPart = strings.ToUpper(elements[1])
			dimValsPart = elements[2]
		} else {
			return "", fmt.Errorf("Failed to parse dimension specs format: %s", formatPart)
		}

		if len(dimPart) == 0 || len(dimValsPart) == 0 {
			return "", fmt.Errorf("Failed to parse dimension specs format: %s", formatPart)
		}

		dimID, err := strconv.ParseInt(dimPart, 10, 64)
		if err != nil {
			return "", fmt.Errorf("Failed to parse dimension specs format: %s", formatPart)
		}
		buffer.WriteString(fmt.Sprintf(`%d,"`, dimID))

		// Just check if is valid format.
		if dimValsPart != "*" {
			for _, dimVal := range strings.Split(dimValsPart, ",") {
				_, err := strconv.ParseInt(dimVal, 10, 64)
				if err != nil {
					return "", fmt.Errorf("Failed to parse dimension values ids: %s", formatPart)
				}
			}
		}

		if len(aggregatorPart) > 0 {
			buffer.WriteString(aggregatorPart)
			buffer.WriteString("(")
			buffer.WriteString(dimValsPart)
			buffer.WriteString(")")
		} else {
			buffer.WriteString(dimValsPart)
		}
		buffer.WriteString(`"]`)
	}
	buffer.WriteString("]")

	// Check the format is valid.
	if err := json.Unmarshal([]byte(buffer.Bytes()), &js); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
