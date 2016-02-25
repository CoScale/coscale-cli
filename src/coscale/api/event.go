package api

import "fmt"

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

// EventData describes the event data uploaded to api
type EventData struct {
	ID         int64
	Timestamp  int64
	Stoptime   int64
	Message    string
	Attribute  string
	Subject    string
	Version    int64
	UpdateTime int64
}

func (e EventData) GetId() int64 {
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

//GetEventData will return the eventdata by the event Id and eventdata Id
func (api *Api) GetEventData(eventId, eventdataId int64, eventData *EventData) error {
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/events/%d/data/get/%d/", api.appID, eventId, eventdataId), nil, false, &eventData); err != nil {
		return err
	}
	return nil
}

func (api *Api) InsertEventData(id int64, message, subject, attribute string, timestamp, stopTime int64) (string, error) {
	data := map[string][]string{
		"message":   {message},
		"timestamp": {fmt.Sprintf("%d", timestamp)},
		"subject":   {subject},
		"attribute": {attribute},
	}
	// add stoptime only if is set
	if stopTime != DEFAULT_INT64_VALUE {
		data["stopTime"] = []string{fmt.Sprintf("%d", stopTime)}
	}

	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/events/%d/data/", api.appID, id), data, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

func (api *Api) UpdateEventData(eventId, eventdataId int64, eventData *EventData) (string, error) {
	data := map[string][]string{
		"message":   {eventData.Message},
		"timestamp": {fmt.Sprintf("%d", eventData.Timestamp)},
		"subject":   {eventData.Subject},
		"attribute": {eventData.Attribute},
		"version":   {fmt.Sprintf("%d", eventData.Version)},
	}
	// add stoptime only if is set
	if eventData.Stoptime != DEFAULT_INT64_VALUE {
		data["stopTime"] = []string{fmt.Sprintf("%d", eventData.Stoptime)}
	}
	var result string
	if err := api.makeCall("PUT", fmt.Sprintf("/api/v1/app/%s/events/%d/data/%d/", api.appID, eventId, eventData.ID), data, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

// DeleteEventData is used to delete a data entry for a event
func (api *Api) DeleteEventData(eventId, eventdataId int64) error {
	if err := api.makeCall("DELETE", fmt.Sprintf("/api/v1/app/%s/events/%d/data/%d/", api.appID, eventId, eventdataId), nil, false, nil); err != nil {
		return err
	}
	return nil
}
