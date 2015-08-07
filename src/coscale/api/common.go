package api

import (
	"fmt"
	"strings"
)

//some of the api calls will be common for different objects
//all the common calls are placed in this source file

//GetObjects will get all the objects (json) specified by objectName. eg: all the metrics or all the servers
func (api *Api) GetObjects(objectName string) (string, error) {
	var result string
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/%ss/", api.appID, objectName), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

// GetObject will return the object (json) specified by objectName that have a certain id
func (api *Api) GetObject(objectName string, id int64) (string, error) {
	var result string
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/%ss/%d/", api.appID, objectName, id), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

//GetObjectRef will put in result a reference to a object specified by objectName and that have a certain id
func (api *Api) GetObjectRef(objectName string, id int64, result Object) error {
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/%ss/%d/", api.appID, objectName, id), nil, false, &result); err != nil {
		return err
	}
	return nil
}

//GetObjectByName will return the object (json) specified by objectName and name
func (api *Api) GetObjectByName(objectName string, name string) (string, error) {
	name = strings.Replace(name, " ", "%20", -1)
	var result string
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/%ss/?selectByName=%s", api.appID, objectName, name), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

//GetObejctRefByName will put in result a reference to the oject specified by objectName and name
func (api *Api) GetObejctRefByName(objectName string, name string, result Object) error {
	name = strings.Replace(name, " ", "%20", -1)
	objects := []*Object{&result}
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/%ss/?selectByName=%s", api.appID, objectName, name), nil, false, &objects); err != nil {
		return err
	}
	if len(objects) == 0 {
		return fmt.Errorf("Not Found")
	} else {
		result = *objects[0]
		return nil
	}
}

//DeleteObject will delete a object
func (api *Api) DeleteObject(objectName string, object *Object) error {
	if err := api.makeCall("DELETE", fmt.Sprintf("/api/v1/app/%s/%ss/%d/", api.appID, objectName, (*object).GetId()), nil, false, nil); err != nil {
		return err
	}
	return nil
}

// AddObjectToGroup adds a object (metric, event, etc) to a group of objects.
func (api *Api) AddObjectToGroup(objectName string, object Object, group Object) error {
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/%sgroups/%d/%ss/%d/", api.appID, objectName, group.GetId(), objectName, object.GetId()), nil, false, nil); err != nil {
		if IsRequestError(err) {
			// The object is already in the group. Ignore this error.
		} else {
			return err
		}
	}
	return nil
}

// AddObjectToGroup remove a object (metric, event, etc) from a group of objects.
func (api *Api) DeleteObjectFromGroup(objectName string, object Object, group Object) error {
	return api.makeCall("DELETE", fmt.Sprintf("/api/v1/app/%s/%sgroups/%d/%ss/%d/", api.appID, objectName, group.GetId(), objectName, object.GetId()), nil, false, nil)
}
