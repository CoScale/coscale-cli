package api

import (
	"fmt"
	"strconv"
)

// Server describes the server object on the API.
type Server struct {
	ID          int64
	Name        string
	Description string
	Type        string
	Source      string
	Attributes  []*ServerAttribute
	State       string
	Version     int64
	Children    map[string]*Server
}

// GetId returns the Id of the Server.
func (e Server) GetId() int64 {
	return e.ID
}

// ServerAttribute describes the server attribute object on the API.
type ServerAttribute struct {
	ID     int64
	Key    string
	Value  string
	Source string
}

// ServerGroup describes the server group object on the API.
type ServerGroup struct {
	ID          int64
	Name        string
	Description string
	Type        string
	Source      string
	State       string
	Version     int64
}

// GetId returns the Id of the ServerGroup.
func (e ServerGroup) GetId() int64 {
	return e.ID
}

// CreateServer creates a new Server using the API.
func (api *Api) CreateServer(name string, description string, serverType string) (string, error) {
	data := map[string][]string{
		"name":        {name},
		"description": {description},
		"type":        {serverType},
		"source":      {GetSource()},
	}
	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/servers/", api.AppID), data, true, &result); err != nil {
		if duplicate, id := IsDuplicate(err); duplicate {
			return api.GetObject("server", id)
		}
		return "", err
	}
	return result, nil
}

// UpdateServer updates all fields on an existing Server using the API.
func (api *Api) UpdateServer(server *Server) (string, error) {
	data := map[string][]string{
		"name":        {server.Name},
		"description": {server.Description},
		"type":        {server.Type},
		"source":      {server.Source},
		"state":       {server.State},
		"version":     {strconv.FormatInt(server.Version, 10)},
	}
	var result string
	if err := api.makeCall("PUT", fmt.Sprintf("/api/v1/app/%s/servers/%d/", api.AppID, server.ID), data, true, &result); err != nil {
		return "", err
	}
	return api.GetObject("server", server.ID)
}

// CreateServerGroup creates a new ServerGroup using the API.
func (api *Api) CreateServerGroup(name, description, Type, state string) (string, error) {
	data := map[string][]string{
		"name":        {name},
		"description": {description},
		"type":        {Type},
		"state":       {state},
		"source":      {GetSource()},
	}

	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/servergroups/", api.AppID), data, true, &result); err != nil {
		if duplicate, id := IsDuplicate(err); duplicate {
			return api.GetObject("servergroup", id)
		}
		return "", err
	}
	return result, nil
}

// UpdateServerGroup updates all fields of an existing ServerGroup using the API.
func (api *Api) UpdateServerGroup(serverGroup *ServerGroup) (string, error) {
	data := map[string][]string{
		"name":        {serverGroup.Name},
		"description": {serverGroup.Description},
		"type":        {serverGroup.Type},
		"state":       {serverGroup.State},
		"source":      {serverGroup.Source},
		"version":     {strconv.FormatInt(serverGroup.Version, 10)},
	}
	var result string
	if err := api.makeCall("PUT", fmt.Sprintf("/api/v1/app/%s/servergroups/%d/", api.AppID, serverGroup.ID), data, true, &result); err != nil {
		return "", err
	}
	return api.GetObject("servergroup", serverGroup.ID)
}
