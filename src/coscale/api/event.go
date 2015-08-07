package api

import (
	"fmt"
	"time"
)

type Event struct {
	ID                    int64
	Name                  string
	Description           string
	AttributeDescriptions string
	Type                  string
	Source                string
	State                 string
	Version               int64
}

func (e Event) GetId() int64 {
	return e.ID
}

func (api *Api) CreateEvent(name, description, attributeDescriptions, source, typeString string) (string, error) {
	data := map[string][]string{
		"name":                  {name},
		"description":           {description},
		"attributeDescriptions": {attributeDescriptions},
		"type":                  {typeString},
		"source":                {GetSource()},
	}
	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/events/", api.appID), data, true, &result); err != nil {
		if duplicate, id := IsDuplicate(err); duplicate {
			return api.GetObject("event", id)
		}
		return "", err
	}
	return result, nil
}

func (api *Api) UpdateEvent(event *Event) (string, error) {
	data := map[string][]string{
		"name":                  {event.Name},
		"description":           {event.Description},
		"attributeDescriptions": {event.AttributeDescriptions},
		"type":                  {event.Type},
		"source":                {event.Source},
		"version":               {fmt.Sprintf("%d", event.Version)},
	}
	var result string
	if err := api.makeCall("PUT", fmt.Sprintf("/api/v1/app/%s/events/%d/", api.appID, event.ID), data, true, &result); err != nil {
		return "", err
	}
	return api.GetObject("event", event.ID)
}

func (api *Api) DeleteEvent(event *Event) error {
	if err := api.makeCall("DELETE", fmt.Sprintf("/api/v1/app/%s/events/%d/", api.appID, event.ID), nil, false, nil); err != nil {
		return err
	}
	return nil
}

func (api *Api) InsertEventData(id int64, message, subject, attribute string, timestamp, stopTime int64) (string, error) {
	now := int64(time.Now().Unix())
	data := map[string][]string{
		"message":   {message},
		"timestamp": {fmt.Sprintf("%d", timestamp - now)},
		"stopTime":  {fmt.Sprintf("%d", stopTime - now)},
		"subject":   {subject},
		"attribute": {attribute},
	}
	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/events/%d/data/", api.appID, id), data, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

func (api *Api) GetEventData(id int64) (string, error) {
	var result string
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/events/%d/data/?start=1434431175", api.appID, id), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}
