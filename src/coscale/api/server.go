package api

import (
	"fmt"
	"strconv"
	"strings"
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
	ParentID    int64
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

// GetServerGroupByPath returns a server group by the hierarchy.
func (api *Api) GetServerGroupByPath(path string) (string, error) {

	var serverGroup *ServerGroup

	groupNames := strings.Split(path, "/")
	for i, groupName := range groupNames {

		var query string

		if i == 0 {
			query = "selectByRoot=true"
		} else if serverGroup != nil && serverGroup.ID > 0 {
			query = fmt.Sprintf("selectByParent_id=%d", serverGroup.ID)
		} else {
			return "[]", nil
		}

		api.SetQueryString(query)
		result, err := api.GetObjectByName("servergroup", groupName)
		if err != nil {
			return "", err
		}
		// For the last group just return it.
		if i == len(groupNames)-1 {
			return result, nil
		}

		tmp := make([]*ServerGroup, 1)
		err = api.HandleResponse([]byte(result), false, &tmp)
		if err != nil {
			return "", err
		}
		if len(tmp) != 1 {
			return "[]", nil
		}
		serverGroup = tmp[0]
	}

	return "[]", nil
}

// CreateServerGroup creates a new ServerGroup using the API.
func (api *Api) CreateServerGroup(name, description, Type, state string, parentID int64) (string, error) {
	data := map[string][]string{
		"name":        {name},
		"description": {description},
		"type":        {Type},
		"state":       {state},
		"source":      {GetSource()},
	}

	// Set the parentId value if its provided.
	if parentID != -1 {
		data["parentId"] = []string{fmt.Sprintf("%d", parentID)}
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
		"parentId":    {fmt.Sprintf("%d", serverGroup.ParentID)},
		"version":     {strconv.FormatInt(serverGroup.Version, 10)},
	}
	var result string
	if err := api.makeCall("PUT", fmt.Sprintf("/api/v1/app/%s/servergroups/%d/", api.AppID, serverGroup.ID), data, true, &result); err != nil {
		return "", err
	}
	return api.GetObject("servergroup", serverGroup.ID)
}
