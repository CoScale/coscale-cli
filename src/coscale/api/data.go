package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// dataPattern matches multiple datapoints because of that it ends with ';'.
var dataPattern = regexp.MustCompile(`(M[0-9]+):([AS]{1}[0-9]*):(.+?)(?:\:(\{(?:.*?)\}))?;`)

// matchDouble matches the data points for double values. ends with ',' because we can have multiple datapoints.
var matchDouble = regexp.MustCompile(`(-?[0-9]+):([0-9.]+),`).FindAllStringSubmatch

// matchHistogram matches data points for histogram, ends with ',' just to be consistent with doubleValuePattern.
var matchHistogram = regexp.MustCompile(`(-?[0-9]+):(\[[0-9]+,[0-9]+,\[[0-9.,]+\]\]),`).FindAllStringSubmatch

// ApiData contains the required fields for a data insert on the API.
type ApiData struct {
	MetricID        int64
	SubjectID       string
	Data            []DataPoint
	DimensionValues map[string]string
}

// HasDimensions will check the ApiData has exactly those dimensions.
func (data *ApiData) HasDimensions(dimensions map[string]string) bool {
	if len(data.DimensionValues) != len(dimensions) {
		return false
	}
	// Check if the two maps are equal.
	for key, val := range data.DimensionValues {
		if dim, ok := dimensions[key]; !(ok && val == dim) {
			return false
		}
	}
	return true
}

// DataPoint is a single data point for a data insert on the API.
type DataPoint struct {
	SecondsAgo int
	Data       string
}

// String formats the DataPoint into the format required by the API.
func (data *DataPoint) String() string {
	return fmt.Sprintf("[%d,%s]", data.SecondsAgo, data.Data)
}

// String formats the ApiData into the format required by the API.
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

// splitPoints checks the format of the dataPoint and splits the string into multiple data points.
func splitPoints(dataPoints string, timeInSecAgo bool) (map[string][]*ApiData, error) {
	if len(dataPoints) == 0 {
		return nil, fmt.Errorf("Bad datapoint format")
	}

	// add semicolon at the end if is necessary, it will help for better matching.
	if !strings.HasSuffix(dataPoints, ";") {
		dataPoints += `;`
	}

	// Match the received data against the dataPattern and extract the expected fields.
	points := dataPattern.FindAllStringSubmatch(dataPoints, -1)

	if len(points) == 0 {
		return nil, fmt.Errorf("Bad datapoint format")
	}

	data := make(map[string][]*ApiData)
	for _, point := range points {

		// The point should have a certain length even if dimension values are missing.
		if len(point) != 5 {
			return nil, fmt.Errorf("Bad datapoint format %s", point)
		}

		// Parse the metric id into int64
		metricID, err := strconv.ParseInt(point[1][1:], 10, 64)
		if err != nil {
			return nil, err
		}

		subjectID := point[2]

		valueStr := point[3]

		// Should not happen because of regex
		if len(valueStr) == 0 {
			return nil, fmt.Errorf("Bad datapoint value is empty")
		}

		// Remove the brackets if is the case.
		if strings.HasPrefix(valueStr, "[") {
			valueStr = strings.TrimPrefix(valueStr, "[")
			valueStr = strings.TrimSuffix(valueStr, "]")
		}

		// add colon at the end if is necessary, it will help for better matching.
		if !strings.HasSuffix(valueStr, ",") {
			valueStr += `,`
		}

		var dimValues map[string]string
		if len(point[4]) > 0 {
			if err := json.Unmarshal([]byte(point[4]), &dimValues); err != nil {
				return nil, err
			}
		}

		// Match the data points.
		var values [][]string

		values = matchDouble(valueStr, -1)

		if len(values) == 0 {
			values = matchHistogram(valueStr, -1)
		}

		if len(values) == 0 {
			return nil, fmt.Errorf("Bad datapoint format at %s, on matching datapoint type", point[3])
		}

		// create the new ApiData object.
		newAPIData := &ApiData{metricID, subjectID, []DataPoint{}, dimValues}
		for _, value := range values {

			if len(value) != 3 {
				return nil, fmt.Errorf("Bad datapoint format at %s", value)
			}

			time, err := strconv.Atoi(value[1])
			if err != nil {
				return nil, fmt.Errorf("Bad datapoint format while parsing the timestamp: %s", value[1])
			}

			// Convert the time format if timeInSecAgo is true. (Is for the deprecated call)
			if timeInSecAgo {
				time = -time
			}

			// Save the parsed DataPoint.
			newAPIData.Data = append(newAPIData.Data, DataPoint{time, value[2]})
		}

		// find the right place for this ApiData in the result
		// if a callData for this subjectId exists, then the new data belongs to it.
		if subjectData, found := data[subjectID]; found {
			// search for an existing ApiData for this metricId
			var foundMetricID bool
			for _, apiData := range subjectData {
				if apiData.MetricID == metricID && apiData.HasDimensions(dimValues) {
					// found the apiData, append to it just the new dataPoint
					apiData.Data = append(apiData.Data, newAPIData.Data...)
					foundMetricID = true
					break
				}
			}
			// a apiData with the same metric Id doesn't exists, create a new one.
			if !foundMetricID {
				data[subjectID] = append(data[subjectID], newAPIData)
			}
		} else {
			// this is a new subjectId, create a new callData.
			data[subjectID] = []*ApiData{newAPIData}
		}
	}

	return data, nil
}

// ParseDataPoint will parse a dataPoints which is provided by user on the command line
// the format is this:
// <METRIC>:<SUBJECT>:<TIME>:[<SAMPLES>,<PERCENTILE WIDTH>,[<PERCENTILE DATA>]]:<{"DIMENSIONS": "JSON"}>
// eg: M1:S1:-60:[100,50,[1,2,3,4,5,6]]:{"Queue":"q1","Data Center":"data center 1"}
// Multiple dataPoints can be splited by semicolons
// eg: M1:S1:-60:1.3:{"Queue":"q1","Data Center":"data center 1"};M2:S1:[-60:1.2,0:1.1]
// If timeInSecAgo is true, the time should be positive and is the number of seconds ago. Otherwise
// it is the time format as defined by the api.
func ParseDataPoint(dataPoints string, timeInSecAgo bool) ([]map[string][]*ApiData, error) {
	data, err := splitPoints(dataPoints, timeInSecAgo)
	if err != nil {
		return nil, err
	}

	_, uncompressedSize, err := serializeAPIData(data)
	if err != nil {
		return nil, err
	}

	// If the size is not exceed, we will have only one batch.
	if uncompressedSize < MaxUploadSize {
		return []map[string][]*ApiData{data}, nil
	}
	return getBatchedAPIData(data)
}

// getBatchedAPIData returns batches of data so that each doesn't exceed uploadSize
// it returns the last error for exceeding upload size
func getBatchedAPIData(data map[string][]*ApiData) ([]map[string][]*ApiData, error) {
	var batches []map[string][]*ApiData
	var err error

	uploadSize := int(0.80 * float64(MaxUploadSize))

	batchSize := 0
	batchCounter := 0
	if len(data) > 0 {
		// we make the current batch
		batches = append(batches, make(map[string][]*ApiData))
	}
	for key, values := range data {
		lastIndex := -1 // index of last added value for current key
		for index, value := range values {
			batchSize += len(value.String()) + len(key) // batch size if this value were to be added
			if batchSize > uploadSize {
				// save current batch
				if lastIndex == -1 && index == 0 {
					// if this is the first time we should add data, but actually there is nothing to add
				} else {
					if lastIndex == -1 {
						lastIndex = 0
					}
					if batchCounter > len(batches)-1 {
						batches = append(batches, make(map[string][]*ApiData))
					}
					batches[batchCounter][key] = values[lastIndex:index]
				}
				// move to next batch
				batchCounter++
				lastIndex = index
				batchSize = len(value.String()) // current value length
				if batchSize > uploadSize {     // this means current value is too big, unlikely to happen, store and skip
					err = fmt.Errorf("Metric data too big, upload size: %d, server id: %s, value: %v", uploadSize, key, value)
					batchSize = 0
					lastIndex = index + 1
				}
			}
		}
		// we add to current batch remaining values
		if lastIndex == -1 { // if data was never added
			lastIndex = 0
		}
		if lastIndex < len(values) {
			if batchCounter > len(batches)-1 {
				batches = append(batches, make(map[string][]*ApiData))
			}
			batches[batchCounter][key] = values[lastIndex:len(values)]
			batchSize += len(values[len(values)-1].String())
		}
	}
	return batches, err
}

// PassThru is a wrapper for gzip.
type PassThru struct {
	io.WriteCloser
	total int // Total  bytes transferred
}

// Write will append the content.
func (pt *PassThru) Write(p []byte) (int, error) {
	n, err := pt.WriteCloser.Write(p)
	pt.total += len(p)
	return n, err
}

// Close the writer.
func (pt *PassThru) Close() error {
	return pt.WriteCloser.Close()
}

// serializeAPIData will base64 and gzip the api data for a single call.
func serializeAPIData(data map[string][]*ApiData) (string, int, error) {
	var buffer bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buffer)
	gzipWriter := &PassThru{gzip.NewWriter(encoder), 0}

	gzipWriter.Write([]byte("["))

	counter := 0
	for _, childData := range data {
		if len(childData) == 0 {
			continue
		}

		if counter > 0 {
			gzipWriter.Write([]byte(","))
		}

		for i, apiData := range childData {
			if i > 0 {
				gzipWriter.Write([]byte(","))
			}
			gzipWriter.Write([]byte(apiData.String()))
		}

		counter++
	}
	gzipWriter.Write([]byte("]"))

	errGzip := gzipWriter.Close()
	errEnc := encoder.Close()

	if errGzip != nil {
		return "", 0, errGzip
	}
	if errEnc != nil {
		return "", 0, errEnc
	}

	return buffer.String(), gzipWriter.total, nil
}

// InsertData inserts a batch of data. Returns the number of pending actions.
func (api *Api) InsertData(data map[string][]*ApiData) (string, error) {

	serializedData, _, err := serializeAPIData(data)
	if err != nil {
		return "", err
	}

	postData := map[string][]string{
		"cdata": {serializedData},
	}
	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/data/", api.AppID), postData, true, &result); err != nil {
		return "", err
	}

	return result, nil
}

// GetData performs an API call to retrieve data from the API.
func (api *Api) GetData(start, stop int, metricId int64, subjectIds, aggregator, viewType, dimensionsSpecs string, aggregateSubjects bool) (string, error) {
	postData := map[string][]string{
		"data": {getBatchData(start, stop, metricId, subjectIds, aggregator, viewType, dimensionsSpecs, aggregateSubjects)},
	}
	var result string
	if err := api.makeCall("POST", fmt.Sprintf("/api/v1/app/%s/data/dimension/getCalculated/", api.AppID), postData, true, &result); err != nil {
		return "", err
	}
	return result, nil
}

// getBatchData make the json object(with the informations provided on command line) required for GetData-getBatch request
func getBatchData(start, stop int, metricId int64, subjectIds, aggregator, viewType, dimensionsSpecs string, aggregateSubjects bool) string {
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
	buffer.WriteString(fmt.Sprintf(`", "aggregator":"%s", "viewtype":"%s", "dimensionsSpecs":%s, "aggregateSubjects":%t}]}`, aggregator, viewType, dimensionsSpecs, aggregateSubjects))
	return buffer.String()
}
