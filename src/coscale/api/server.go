package api

import (
	"fmt"
	"strconv"
)

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

func (e Server) GetId() int64 {
	return e.ID
}

type ServerAttribute struct {
	ID     int64
	Key    string
	Value  string
	Source string
}

type ServerGroup struct {
	ID          int64
	Name        string
	Description string
	Type        string
	Source      string
	State       string
	Version     int64
}

func (e ServerGroup) GetId() int64 {
	return e.ID
}

// CreateServer creates a new server.
func (api *Api) CreateServer(name string, description string, serverType, source string) (string, error) {
	data := map[string][]string{
		"name":        {name},
		"description": {description},
		"type":        {serverType},
		"source":      {source},
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

// CreateServerGroup creates a new server group.
func (api *Api) CreateServerGroup(name, description, Type, state, source string) (string, error) {
	data := map[string][]string{
		"name":        {name},
		"description": {description},
		"type":        {Type},
		"state":       {state},
		"source":      {source},
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

// UpdateServerGroup updates the name of the servergroup.
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
