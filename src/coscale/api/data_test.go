package api

import (
	"reflect"
	"testing"
)

// Combine histogram with double data. For double data we insert for the same metric id in both methods:
// - as another data point splitted by ',' and
// - in the old way splited by ';' and providing metric id and subject.
var data = `M2:S2:-60:1.1,0:1.2;M2:S2:[-120:9.1];M1:S1:-60:[100,50,[1,2,3,4,5,6]]:{"Queue":"q1","Data Center":"data center 1"};M1:S1:[-120:[100,50,[1,2,3,4,5,7]],-180:[100,50,[1,2,3,4,5,8]]]:{"Queue":"q1","Data Center":"data center 1"}`

var expected = []map[string][]*ApiData{{
	"S1": []*ApiData{
		&ApiData{
			MetricID:  1,
			SubjectID: "S1",
			Data: []DataPoint{
				{
					-60,
					"[100,50,[1,2,3,4,5,6]]",
				},
				{
					-120,
					"[100,50,[1,2,3,4,5,7]]",
				},
				{
					-180,
					"[100,50,[1,2,3,4,5,8]]",
				},
			},
			DimensionValues: map[string]string{
				"Queue":       "q1",
				"Data Center": "data center 1",
			},
		},
	},
	"S2": []*ApiData{
		&ApiData{
			MetricID:  2,
			SubjectID: "S2",
			Data: []DataPoint{
				{
					-60,
					"1.1",
				},
				{
					-0,
					"1.2",
				},
				{
					-120,
					"9.1",
				},
			},
			DimensionValues: nil,
		},
	},
},
}

// Test ParseDataPoint.
func TestParseDataPoint(t *testing.T) {

	// Correct case.
	obtained, err := ParseDataPoint(data, false)
	if err != nil {
		t.Fatalf("Error occured while parsing data: %s", err)
	}
	if !reflect.DeepEqual(expected, obtained) {
		t.Fatalf("expected: \n%v\n, found: \n%v\n", expected, obtained)
	}

	// Bad format.
	obtained2, err := ParseDataPoint(`M2:S2::{"Queue":"q1","Data Center":"data center 1"}`, false)
	if err == nil {
		t.Fatalf("Expected error.")
	}
	if obtained2 != nil {
		t.Fatalf("expected: \n%v\n, found: \n%v\n", "[]", obtained2)
	}

	obtained3, err := ParseDataPoint(`M2:S2:32:{"Queue":"q1","Data Center":"data center 1"}`, false)
	if err == nil {
		t.Fatalf("Expected error.")
	}
	if obtained3 != nil {
		t.Fatalf("expected: \n%v\n, found: \n%v\n", "[]", obtained3)
	}
}
