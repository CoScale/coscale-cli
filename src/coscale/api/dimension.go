package api

import "fmt"

// Dimension is a metric dimension.
type Dimension struct {
	ID   int64
	Name string
	// Link
}

func (d *Dimension) String() string {
	return fmt.Sprintf("%#v", d)
}

// GetDimension gets one dimension by its id.
func (api *Api) GetDimension(id int64) (string, error) {
	var result string
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/dimensions/%d/", api.AppID, id), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

// GetDimensions gets the dimensions for a metric
func (api *Api) GetDimensions(metricID int64) (string, error) {
	var result string
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/metrics/%d/dimensions/", api.AppID, metricID), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

// CreateDimension creates a new dimension.
func (api *Api) CreateDimension(name string) (string, error) {
	data := map[string][]string{
		"name": {name},
	}

	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/dimensions/", api.AppID), data, true, &result); err != nil {
		if duplicate, id := IsDuplicate(err); duplicate {
			return api.GetDimension(id)
		}
		return "", err
	}

	return result, nil
}
