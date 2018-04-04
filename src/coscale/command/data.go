package command

import (
	"coscale/api"
	"os"
)

var dataObjectName = "data"

// DataObject defines the data command on the CLI.
var DataObject = NewCommand(dataObjectName, "data <action> [--<field>='<data>']", DataActions)

// DataActions defines the data actions on the CLI.
var DataActions = []*Command{
	{
		Name:      "get",
		UsageLine: `data get (--id --subjectIds) [--start --stop --aggregator --viewtype --aggregateSubjects]`,
		Long: `
Retrieve a batch of data from the datastore.

The flags for get data action are:
Mandatory:
	--id
		Metric id.
	--subjectIds
		The subject string eg. s1 for server 1, g2 for servergroup 2, a for application.
Optional:
	--start
		The start timestamp in seconds ago(negative values) or unix timestamp (positive values). [default: 0]
	--stop
		The stop timestamp in seconds ago(negative values) or unix timestamp (positive values). [default: 0]
	--aggregator
		The data aggregator(AVG, MIN, MAX) used to specify vertical aggregation of timeseries. [default: AVG]
	--viewtype
		The view type defines how the data will be shown. [default: DEFAULT]
			DEFAULT: interpretation depends on metricType
			AVG: returns average data
			MIN: returns minimal data
			MAX: returns maximal data
			RATE: returns data as a rate (always in #/s)
			COUNT: returns number of occurences since previous timestamp
	--dimensionsSpecs
		[[<dimension_id>, <dimension_spec>],
		 [<dimension_id>, <dimension_spec>], ...]

 		<dimension_spec> = "*" will return data for all dimensionvalues separately
 					| "<dimension_value_id>, <dimension_value_id>, ..."
					| "<aggregator>(*)"
					| "<aggregator>([<dimension_value_id>, <dimension_value_id>, ...])"

	    aggregators: AVG, MIN, MAX

		e.g.: --dimensionsSpecs='[[1,"AVG(*)"]]'
		      --dimensionsSpecs='[[2,"*"]]'
		      --dimensionsSpecs='[[3,"11,12,13"],[4,"21,22,23"]]'

	--aggregateSubjects
		Boolean that indicates if the aggregated value over all subjectIds should be returned. [default: false]
`,
		Run: func(cmd *Command, args []string) {
			var subjectIds, aggregator, viewType, dimensionsSpecs string
			var aggregateSubjects bool
			var start, stop int
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier for metric.")
			cmd.Flag.IntVar(&start, "start", 0, "The start timestamp in seconds ago.")
			cmd.Flag.IntVar(&stop, "stop", 0, "The stop timestamp in seconds ago.")
			cmd.Flag.StringVar(&subjectIds, "subjectIds", DEFAULT_STRING_FLAG_VALUE, "The subject string.")
			cmd.Flag.StringVar(&aggregator, "aggregator", "AVG", "The data aggregator (AVG, MIN, MAX).")
			cmd.Flag.StringVar(&viewType, "viewType", "DEFAULT", "Defines how the data will be shown.")
			cmd.Flag.StringVar(&dimensionsSpecs, "dimensionsSpecs", "[]", "JSON containing ids of the dimensions.")
			cmd.Flag.BoolVar(&aggregateSubjects, "aggregateSubjects", false, "Boolean that indicates if the aggregated value over all subjectIds should be returned.")
			cmd.ParseArgs(args)
			if subjectIds == DEFAULT_STRING_FLAG_VALUE || id == -1 {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			cmd.PrintResult(cmd.Capi.GetData(start, stop, id, subjectIds, aggregator, viewType, dimensionsSpecs, aggregateSubjects))
		},
	},
	{
		Name:      "insert",
		UsageLine: `data insert (--data)`,
		Long: `
Insert data for metrics into the datastore.

The flags for data insert action are:
Mandatory:
	--data
		To send data for DOUBLE metric data typ use the following format:
			"M<metric id>:S<subject Id>:<time>:<value/s>"
			eg:	--data="M1:S100:1454580954:1.2"

		To send data for HISTOGRAM metric data type use the following format:
			"M<metric id>:S<subject Id>:<seconds ago>:[<no of samples>,<percentile width>,[<percentile data>]]"
			eg: --data="M1:S1:-60:[100,50,[1,2,3,4,5,6]]"

		Sending multiple data points for the same metric and subject is possible using the folowing format:
			--data="M1:S100:[-60:1.2,0:1.1]"

		Sending multiple data entries is possible by using semicolon as separator.
			eg: --data="M1:S100:-60:1.2;M2:S100:0:2"

		The time is formatted as follows:
		    Positive numbers are interpreted as unix timestamps in seconds.
		    Zero is interpreted as the current time.
		    Negative numbers are interpreted as a seconds ago from the current time.
		Metric dimensions enables us to show metrics at different levels. For example for RabbitMQ
			we want to show the total number of queued messages, but we also want to be able
			to split these into the number of queued messages per queue.
			eg: --data='M1:S1:-60:1.3:{"Queue":"q1","Data Center":"data center 1"};M2:S1:-60:1.2'


Deprecated:
	--datapoint
		To send data for DOUBLE metric data type use the following format:
			"M<metric id>:S<subject Id>:<seconds ago>:<value/s>"
			eg:	--datapoint="M1:S100:120:1.2

		To send data for HISTOGRAM metric data type use the following format:
			"M<metric id>:S<subject Id>:<seconds ago>:[<no of samples>,<percentile width>,[<percentile data>]]"
			eg: --datapoint="M1:S1:60:[100,50,[1,2,3,4,5,6]]"
`,
		Run: func(cmd *Command, args []string) {
			var datapoint string
			var data string

			cmd.Flag.Usage = func() { cmd.PrintUsage() }

			cmd.Flag.StringVar(&datapoint, "datapoint", DEFAULT_STRING_FLAG_VALUE, "")
			cmd.Flag.StringVar(&data, "data", DEFAULT_STRING_FLAG_VALUE, "")
			cmd.ParseArgs(args)

			if datapoint == DEFAULT_STRING_FLAG_VALUE && data == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}

			timeInSecAgo := false
			if data == DEFAULT_STRING_FLAG_VALUE {
				timeInSecAgo = true
				data = datapoint
			}

			// ParseDataPoint could return data for multiple calls for same metricIDs with different subjectID
			callsData, err := api.ParseDataPoint(data, timeInSecAgo)
			if err != nil {
				cmd.PrintResult("", err)
			}
			var result string
			var resErr error
			// if datapoint contain multiple subject id, then we will have multiple api calls
			for _, callData := range callsData {
				result, resErr = cmd.Capi.InsertData(callData)
				if resErr != nil {
					break
				}
			}
			cmd.PrintResult(result, resErr)
		},
	},
}
