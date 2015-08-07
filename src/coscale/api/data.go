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
//	M2:S100:60:1000,10,0,10,20,30,40,50,60,70,80,90,95,99,100
//	<SAMPLES>,<PERCENTILE WIDTH>,<PERCENTILE DATA>
func ParseDataPoint(dataPoint string) (*ApiData, error) {
	data := strings.Split(dataPoint, ":")
	if len(data) < 4 {
		return nil, fmt.Errorf("Bad datapoint format")
	}

	metricId, err := strconv.ParseInt(data[0][1:], 10, 64)
	if err != nil {
		return nil, err
	}
	secAgo, err := strconv.Atoi(data[2])
	if err != nil {
		return nil, err
	}
	converted := ApiData{metricId, data[1], []DataPoint{DataPoint{-secAgo, data[3]}}}
	return &converted, nil
}

// InsertData inserts a batch of data. Returns the number of pending actions.
func (api *Api) InsertData(data *ApiData) (string, error) {
	postData := map[string][]string{
		"source": {data.SubjectID},
		"data":   {apiDataToString([]*ApiData{data})},
	}
	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/data/", api.appID), postData, true, &result); err != nil {
		return "", err
	}

	return result, nil
}

//make the json object(with the informations provided on command line) required for GetData-getBatch request
func getBatchData(start, stop int, metricId int64, subjectIds, aggregator string, aggregateSubjects bool) string {
	var buffer bytes.Buffer
	var now = int(time.Now().Unix())
	buffer.WriteString(fmt.Sprintf(`{"start":%d, "stop":%d, "ids":[{"metricId":%d, "subjectIds":[`, now-start, now-stop, metricId))
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
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/data/getBatch/", api.appID), postData, true, &result); err != nil {
		return "", err
	}
	return result, nil
}
