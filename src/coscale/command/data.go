package command

import (
	"coscale/api"
	"os"
)

var dataObjectName = "data"
var DataObject = NewCommand(dataObjectName, "data <action> [--<field>='<data>']", DataActions)
var DataActions = []*Command{
	{
		Name:      "get",
		UsageLine: `data get (--id --subjectIds) [--start --stop --aggregator --aggregateSubjects]`,
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
		The data aggregator (AVG, MIN, MAX). [default: AVG]
	--aggregateSubjects
		Boolean that indicates if the aggregated value over all subjectIds should be returned. [default: false]
`,
		Run: func(cmd *Command, args []string) {
			var subjectIds, aggregator string
			var aggregateSubjects bool
			var start, stop int
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier for metric.")
			cmd.Flag.IntVar(&start, "start", 0, "The start timestamp in seconds ago.")
			cmd.Flag.IntVar(&stop, "stop", 0, "The stop timestamp in seconds ago.")
			cmd.Flag.StringVar(&subjectIds, "subjectIds", DEFAULT_FLAG_VALUE, "The subject string")
			cmd.Flag.StringVar(&aggregator, "aggregator", "AVG", "The data aggregator (AVG, MIN, MAX).")
			cmd.Flag.BoolVar(&aggregateSubjects, "aggregateSubjects", false, "Boolean that indicates if the aggregated value over all subjectIds should be returned.")
			cmd.ParseArgs(args)
			if subjectIds == DEFAULT_FLAG_VALUE || id == -1 {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			cmd.PrintResult(cmd.Capi.GetData(start, stop, id, subjectIds, aggregator, aggregateSubjects))
		},
	},
	{
		Name:      "insert",
		UsageLine: `data insert (--datapoint)`,
		Long: `
Create new EventData for a given event.

The flags for data event action are:
Mandatory:
	--datapoint
		Data format is:"M<metric id>:<subject Id>:<seconds ago>:<value/s>" eg:"M1:S100:120:1,2"
`,
		Run: func(cmd *Command, args []string) {
			var datapoint string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&datapoint, "datapoint", DEFAULT_FLAG_VALUE, "")
			cmd.ParseArgs(args)
			if datapoint == DEFAULT_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			data, err := api.ParseDataPoint(datapoint)
			if err != nil {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			cmd.PrintResult(cmd.Capi.InsertData(data))
		},
	},
}
