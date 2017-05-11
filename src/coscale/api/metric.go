package api

import (
	"fmt"
	"strconv"
)

type Metric struct {
	ID          int64
	Name        string
	Description string
	DataType    string
	Period      int
	Unit        string
	Source      string
	Subject     string
	State       string
	Version     int64
}

func (e Metric) GetId() int64 {
	return e.ID
}

type MetricGroup struct {
	ID           int64
	Name         string
	Description  string
	Type         string
	MetricGroups []*MetricGroup
	Source       string
	State        string
	Subject      string
	Version      int64
}

func (e MetricGroup) GetId() int64 {
	return e.ID
}

func (api *Api) CreateMetric(name, description, datatype, unit, subject string, period int) (string, error) {
	data := map[string][]string{
		"name":        {name},
		"description": {description},
		"dataType":    {datatype},
		"period":      {strconv.Itoa(period)},
		"unit":        {unit},
		"subject":     {subject},
		"source":      {GetSource()},
	}
	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/metrics/", api.AppID), data, true, &result); err != nil {
		if duplicate, id := IsDuplicate(err); duplicate {
			return api.GetObject("metric", id)
		}
		return "", err
	}
	return result, nil
}

func (api *Api) UpdateMetric(metric *Metric) (string, error) {
	data := map[string][]string{
		"name":        {metric.Name},
		"description": {metric.Description},
		"dataType":    {metric.DataType},
		"period":      {strconv.Itoa(metric.Period)},
		"unit":        {metric.Unit},
		"subject":     {metric.Subject},
		"source":      {metric.Source},
		"version":     {fmt.Sprintf("%d", metric.Version)},
	}
	var result string
	if err := api.makeCall("PUT", fmt.Sprintf("/api/v1/app/%s/metrics/%d/", api.AppID, metric.ID), data, true, &result); err != nil {
		return "", err
	}
	return api.GetObject("metric", metric.ID)
}

// CreateMetricGroup creates a new metric group.
func (api *Api) CreateMetricGroup(name, description, Type, state, subject string) (string, error) {
	data := map[string][]string{
		"name":        {name},
		"description": {description},
		"type":        {Type},
		"state":       {state},
		"subject":     {subject},
		"source":      {GetSource()},
	}
	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/metricgroups/", api.AppID), data, true, &result); err != nil {
		if duplicate, id := IsDuplicate(err); duplicate {
			return api.GetObject("metricgroup", id)
		}
		return "", err
	}
	return result, nil
}

// UpdateMetricGroup updates the name of the metricgroup.
func (api *Api) UpdateMetricGroup(metricGroup *MetricGroup) (string, error) {
	data := map[string][]string{
		"name":        {metricGroup.Name},
		"description": {metricGroup.Description},
		"type":        {metricGroup.Type},
		"state":       {metricGroup.State},
		"subject":     {metricGroup.Subject},
		"source":      {metricGroup.Source},
		"version":     {fmt.Sprintf("%d", metricGroup.Version)},
	}
	var result string
	if err := api.makeCall("PUT", fmt.Sprintf("/api/v1/app/%s/metricgroups/%d/", api.AppID, metricGroup.ID), data, true, &result); err != nil {
		return "", err
	}
	return api.GetObject("metricgroup", metricGroup.ID)
}

//GetMetricsByGroup will return all the metrics from a metricgroup
func (api *Api) GetMetricsByGroup(metricGroup *MetricGroup) (string, error) {
	var result string
	if err := api.makeCall("GET", fmt.Sprintf("/api/v1/app/%s/metricgroups/%d/metrics/", api.AppID, metricGroup.GetId()), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

// AddMetricDimension adds a dimension to a metric
func (api *Api) AddMetricDimension(metricID, dimensionID int64) (string, error) {
	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/metrics/%d/dimensions/%d/", api.AppID, metricID, dimensionID), nil, true, &result); err != nil {
		return "", err
	}
	return result, nil
}
