package api

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ApiData struct {
	MetricID  int64
	SubjectID string
	Data      []DataPoint
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
	buffer.WriteString("]}")
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

//dataPoint is provided by user on the command line
//the format is this:
// <METRIC>:<SUBJECT>:<TIME>[<SAMPLES>,<PERCENTILE WIDTH>,[<PERCENTILE DATA>]]
// eg: M1:s1:-60:[100,50,[1,2,3,4,5,6]]
// Multiple dataPoint can be splited by semicolons
// eg: M1:s1:-120:0.1;M1:s1:-60:90.9
// If timeInSecAgo is true, the time should be positive and is the number of seconds ago. Otherwise
// it is the time format as defined by the api.
func ParseDataPoint(dataPoint string, timeInSecAgo bool) (map[string][]*ApiData, error) {
	dataPoint = strings.TrimRight(dataPoint, ";")
	// split the data in multiple calls if multiple subjectIds are provided
	// we do this because a call cannot have the same metricId two times, but we should allow adding data
	// for the same metricId on multiple subject ids
	callsData := make(map[string][]*ApiData)

	// the data entries are separated by ";"
	for _, entry := range strings.Split(dataPoint, ";") {
		// parse each data entry
		data := strings.Split(entry, ":")
		if len(data) < 4 {
			return nil, fmt.Errorf("Bad datapoint format")
		}
		metricId, err := strconv.ParseInt(data[0][1:], 10, 64)
		if err != nil {
			return nil, err
		}
		time, err := strconv.Atoi(data[2])
		if err != nil {
			return nil, err
		}
		subjectId := data[1]

		// Convert the time format if timeInSecAgo is true.
		if timeInSecAgo {
			time = -time
		}

		// create the new ApiData object
		newDataPoint := DataPoint{time, data[3]}
		newApiData := &ApiData{metricId, subjectId, []DataPoint{newDataPoint}}

		// find the right place for this ApiData in the result
		// if a callData for this subjectId exists, then the new data belongs to it
		if callData, found := callsData[subjectId]; found {
			// search for an existing ApiData for this metricId
			var foundMetricId bool
			for _, apiData := range callData {
				if apiData.MetricID == metricId {
					// found the apiData, append to it just the new dataPoint
					apiData.Data = append(apiData.Data, newDataPoint)
					foundMetricId = true
					break
				}
			}
			// a apiData with the same metric Id doesn't exists, create a new one
			if !foundMetricId {
				callsData[subjectId] = append(callsData[subjectId], newApiData)
			}
		} else {
			// this is a new subjectId, create a new callData
			callsData[subjectId] = []*ApiData{newApiData}
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
func getBatchData(start, stop int, metricId int64, subjectIds, aggregator string, aggregateSubjects bool) string {
	var buffer bytes.Buffer
	var now = int(time.Now().Unix())
	// negative and null values are seconds ago
	if start <= 0 {
		start = now + start
	}
	if stop <= 0 {
		stop = now + stop
	}
	buffer.WriteString(fmt.Sprintf(`{"start":%d, "stop":%d, "ids":[{"metricId":%d, "subjectIds":[`, start, stop, metricId))
	for i, id := range strings.Split(subjectIds, ",") {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(fmt.Sprintf(`"%s"`, id))
	}
	buffer.WriteString(fmt.Sprintf(`], "aggregator":"%s", "aggregateSubjects":%t}]}`, aggregator, aggregateSubjects))
	return buffer.String()
}

func (api *Api) GetData(start, stop int, metricId int64, subjectIds, aggregator string, aggregateSubjects bool) (string, error) {
	postData := map[string][]string{
		"data": {getBatchData(start, stop, metricId, subjectIds, aggregator, aggregateSubjects)},
	}
	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/data/getBatch/", api.AppID), postData, true, &result); err != nil {
		return "", err
	}
	return result, nil
}
