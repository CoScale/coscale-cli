package api

import (
	"fmt"
	"strings"
)

//some of the api calls will be common for different objects
//all the common calls are placed in this source file

// GetObjects will get all the objects (json) specified by objectName. eg: all the metrics or all the servers
func (api *Api) GetObjects(objectName string) (string, error) {
	var result string
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/%ss/", api.AppID, objectName), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

// GetObject will return the object (json) specified by objectName that have a certain id
func (api *Api) GetObject(objectName string, id int64) (string, error) {
	var result string
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/%ss/%d/", api.AppID, objectName, id), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

// GetObjectFromGroup will return the object (json) specified by objectName from objectGroup that have a certain id
func (api *Api) GetObjectFromGroup(objectGroup, objectName string, groupID, objectID int64) (string, error) {
	var result string
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/%ss/%d/%ss/%d/", api.AppID, objectGroup, groupID, objectName, objectID), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

// GetObjectRef will put in result a reference to a object specified by objectName and that have a certain id
func (api *Api) GetObjectRef(objectName string, id int64, result Object) error {
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/%ss/%d/", api.AppID, objectName, id), nil, false, &result); err != nil {
		return err
	}
	return nil
}

// GetObjectRefFromGroup will return the object specified by objectName from objectGroup that have a certain id
func (api *Api) GetObjectRefFromGroup(objectGroup, objectName string, groupID, objectID int64, result Object) error {
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/%ss/%d/%ss/%d/", api.AppID, objectGroup, groupID, objectName, objectID), nil, false, &result); err != nil {
		return err
	}
	return nil
}

// GetObjectByName will return the object (json) specified by objectName and name
func (api *Api) GetObjectByName(objectName string, name string) (string, error) {
	// In go %% is % escaped, we need to escape the name to work with string fmt.
	name = strings.Replace(name, "%", "%%", -1)
	name = strings.Replace(name, " ", "%20", -1)
	var result string
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/%ss/?selectByName=%s", api.AppID, objectName, name), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

// GetObejctRefByName will put in result a reference to the oject specified by objectName and name
func (api *Api) GetObejctRefByName(objectName string, name string, result Object) error {
	// In go %% is % escaped, we need to escape the name to work with string fmt.
	name = strings.Replace(name, "%", "%%", -1)
	name = strings.Replace(name, " ", "%20", -1)
	objects := []*Object{&result}
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/%ss/?selectByName=%s", api.AppID, objectName, name), nil, false, &objects); err != nil {
		return err
	}
	if len(objects) == 0 {
		return fmt.Errorf("Not Found")
	}
	result = *objects[0]
	return nil
}

// GetObejctRefByNameFromGroup will return the object specified by objectName from objectGroup that have a certain name
func (api *Api) GetObejctRefByNameFromGroup(objectGroup, objectName string, groupID int64, name string, result Object) error {
	// In go %% is % escaped, we need to escape the name to work with string fmt.
	name = strings.Replace(name, "%", "%%", -1)
	name = strings.Replace(name, " ", "%20", -1)
	objects := []*Object{&result}
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/%ss/%d/%ss/?selectByName=%s", api.AppID, objectGroup, groupID, objectName, name), nil, false, &objects); err != nil {
		return err
	}
	if len(objects) == 0 {
		return fmt.Errorf("Not Found")
	}
	result = *objects[0]
	return nil
}

// DeleteObject will delete a object
func (api *Api) DeleteObject(objectName string, object *Object) (string, error) {
	var result string

	if err := api.makeCall("DELETE", fmt.Sprintf("/api/v1/app/%s/%ss/%d/", api.AppID, objectName, (*object).GetId()), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

// AddObjectToGroup adds a object (metric, event, etc) to a group of objects.
func (api *Api) AddObjectToGroup(objectName string, object Object, group Object) (string, error) {
	var objectGroupName = GetObjectGroupName(objectName)

	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/%ss/%d/%ss/%d/", api.AppID, objectGroupName, group.GetId(), objectName, object.GetId()), nil, true, &result); err != nil {
		if IsRequestError(err) {
			// The object is already in the group. Ignore this error.
		} else {
			return "", err
		}
	}
	return result, nil
}

// DeleteObjectFromGroup remove a object (metric, event, etc) from a group of objects.
func (api *Api) DeleteObjectFromGroup(objectName string, object Object, group Object) (string, error) {
	var objectGroupName = GetObjectGroupName(objectName)

	var result string
	if err := api.makeCall("DELETE", fmt.Sprintf("/api/v1/app/%s/%ss/%d/%ss/%d/", api.AppID, objectGroupName, group.GetId(), objectName, object.GetId()), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

// DeleteObjectFromGroupByID remove a object (metric, event, etc) from a group of objects.
func (api *Api) DeleteObjectFromGroupByID(groupName, objectName string, groupID, id int64) (string, error) {
	var result string
	if err := api.makeCall("DELETE", fmt.Sprintf("/api/v1/app/%s/%ss/%d/%ss/%d/", api.AppID, groupName, groupID, objectName, id), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}
