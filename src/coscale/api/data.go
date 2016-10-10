package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// dataPattern is used to split into groups the data inserted by the user.
var dataPattern = regexp.MustCompile(`(M[0-9]+):([ASG]{1}[0-9]*):(-?[0-9]+):([0-9.]+)(?:\:(\{(?:.*?)\}))?;`)

type ApiData struct {
	MetricID        int64
	SubjectID       string
	Data            []DataPoint
	DimensionValues map[string]string
}

// HasDimensions will check the ApiData has exactly those dimensions.
func (a *ApiData) HasDimensions(dimensions map[string]string) bool {
	if len(a.DimensionValues) != len(dimensions) {
		return false
	}
	// Check if the two maps are equal.
	for key, val := range a.DimensionValues {
		if dim, ok := dimensions[key]; !(ok && val == dim) {
			return false
		}
	}
	return true
}

type DataPoint struct {
	SecondsAgo int
	Data       string
}

func (data *DataPoint) String() string {
	return fmt.Sprintf("[%d,%s]", data.SecondsAgo, data.Data)
}

func (data *ApiData) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf(`{"m":%d, "s":"%s", "d":[`, data.MetricID, data.SubjectID))
	for i, d := range data.Data {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(d.String())
	}
	buffer.WriteString("]")

	if len(data.DimensionValues) > 0 {
		buffer.WriteString(`,"dv":{`)
		index := 0
		for dimension, dimensionValue := range data.DimensionValues {
			if index > 0 {
				buffer.WriteString(",")
			}
			buffer.WriteString(fmt.Sprintf(`"%s":"%s"`, dimension, dimensionValue))
			index++
		}
		buffer.WriteString("}")
	}

	buffer.WriteString("}")
	return buffer.String()
}

func apiDataToString(data []*ApiData) string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	for i, d := range data {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(d.String())
	}
	buffer.WriteString("]")
	return buffer.String()
}

// ParseDataPoint will parse a dataPoint which is provided by user on the command line
// the format is this:
// <METRIC>:<SUBJECT>:<TIME>:[<SAMPLES>,<PERCENTILE WIDTH>,[<PERCENTILE DATA>]]:<{"DIMENSIONS": "JSON"}>
// eg: M1:S1:-60:[100,50,[1,2,3,4,5,6]]:{"Queue":"q1","Data Center":"data center 1"}
// Multiple dataPoint can be splited by semicolons
// eg: M1:S1:-60:1.3:{"Queue":"q1","Data Center":"data center 1"};M2:S1:-60:1.2
// If timeInSecAgo is true, the time should be positive and is the number of seconds ago. Otherwise
// it is the time format as defined by the api.
func ParseDataPoint(dataPoint string, timeInSecAgo bool) (map[string][]*ApiData, error) {

	// add semicolon at the end if is neccessary, it will help for better matching.
	if (dataPoint[len(dataPoint)-1]) != ';' {
		dataPoint += `;`
	}

	callsData := make(map[string][]*ApiData)
	// Match the received data against the dataPattern and extract the expected fields.
	matches := dataPattern.FindAllStringSubmatch(dataPoint, -1)

	if len(matches) == 0 {
		return nil, fmt.Errorf("Bad datapoint format")
	}

	for _, match := range matches {

		// The match should have a certain lenght even if dimension values are missing.
		if len(match) != 6 {
			return nil, fmt.Errorf("Bad datapoint format")
		}

		// Parse the metric id into int64
		metricID, err := strconv.ParseInt(match[1][1:], 10, 64)
		if err != nil {
			return nil, err
		}

		subjectID := match[2]

		time, err := strconv.Atoi(match[3])
		if err != nil {
			return nil, err
		}

		// Convert the time format if timeInSecAgo is true. (Is for the deprecated call)
		if timeInSecAgo {
			time = -time
		}
		var dimValues map[string]string
		if len(match[5]) > 0 {
			if err := json.Unmarshal([]byte(match[5]), &dimValues); err != nil {
				return nil, err
			}
		}

		// create the new ApiData object
		newDataPoint := DataPoint{time, match[4]}
		newAPIData := &ApiData{metricID, subjectID, []DataPoint{newDataPoint}, dimValues}

		// find the right place for this ApiData in the result
		// if a callData for this subjectId exists, then the new data belongs to it
		if callData, found := callsData[subjectID]; found {
			// search for an existing ApiData for this metricId
			var foundMetricID bool
			for _, apiData := range callData {
				if apiData.MetricID == metricID && apiData.HasDimensions(dimValues) {
					// found the apiData, append to it just the new dataPoint
					apiData.Data = append(apiData.Data, newDataPoint)
					foundMetricID = true
					break
				}
			}
			// a apiData with the same metric Id doesn't exists, create a new one
			if !foundMetricID {
				callsData[subjectID] = append(callsData[subjectID], newAPIData)
			}
		} else {
			// this is a new subjectId, create a new callData
			callsData[subjectID] = []*ApiData{newAPIData}
		}
	}
	return callsData, nil
}

// InsertData inserts a batch of data. Returns the number of pending actions.
func (api *Api) InsertData(data []*ApiData) (string, error) {
	postData := map[string][]string{
		"data": {apiDataToString(data)},
	}
	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/data/", api.AppID), postData, true, &result); err != nil {
		return "", err
	}

	return result, nil
}

//make the json object(with the informations provided on command line) required for GetData-getBatch request
func getBatchData(start, stop int, metricId int64, subjectIds, aggregator, dimensionsSpecs string, aggregateSubjects bool) string {
	var buffer bytes.Buffer
	var now = int(time.Now().Unix())
	// negative and null values are seconds ago
	if start <= 0 {
		start = now + start
	}
	if stop <= 0 {
		stop = now + stop
	}
	buffer.WriteString(fmt.Sprintf(`{"start":%d, "stop":%d, "ids":[{"metricId":%d, "subjects":"`, start, stop, metricId))
	for i, id := range strings.Split(subjectIds, ",") {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(fmt.Sprintf(`%s`, id))
	}
	buffer.WriteString(fmt.Sprintf(`", "aggregator":"%s", "dimensionsSpecs":%s, "aggregateSubjects":%t}]}`, aggregator, dimensionsSpecs, aggregateSubjects))
	return buffer.String()
}

func (api *Api) GetData(start, stop int, metricId int64, subjectIds, aggregator, dimensionsSpecs string, aggregateSubjects bool) (string, error) {
	postData := map[string][]string{
		"data": {getBatchData(start, stop, metricId, subjectIds, aggregator, dimensionsSpecs, aggregateSubjects)},
	}
	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/data/dimension/getCalculated/", api.AppID), postData, true, &result); err != nil {
		return "", err
	}
	return result, nil
}
