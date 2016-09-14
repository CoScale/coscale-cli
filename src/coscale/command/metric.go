package command

import (
	"coscale/api"
	"os"
)

var MetricObject = NewCommand("metric", "metric <action> [--<field>='<data>']", MetricActions)

var MetricActions = []*Command{
	ListCmd("metric"),
	GetCmd("metric"),
	{
		Name:      "listbygroup",
		UsageLine: `metric listbygroup (--id | --name)`,
		Long: `
Get all metrics from a metric group

The flags for listbygroup metric action are:

Mandatory:
	--id
		Unique identifier for a metricgroup
	or 
	--name 
		specify the name of the metrigroup.
`,
		Run: func(cmd *Command, args []string) {
			var id int64
			var name string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier")
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the metric group.")

			cmd.ParseArgs(args)

			var metricGroupObj = &api.MetricGroup{}
			var err error
			if id != -1 {
				err = cmd.Capi.GetObjectRef("metricgroup", id, metricGroupObj)
			} else if name != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName("metricgroup", name, metricGroupObj)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}

			cmd.PrintResult(cmd.Capi.GetMetricsByGroup(metricGroupObj))
		},
	},
	{
		Name:      "new",
		UsageLine: `metric new (--name --dataType --subject) [--period --description --unit --attachTo --source]`,
		Long: `
Create a new CoScale metric object.

The flags for new metric action are:

Mandatory:
	--name 
		specify the name of the metric.
	--dataType
		The following data types are defined: "LONG", "DOUBLE", "HISTOGRAM".
	--subject
		A metric is defined on either a "SERVER", "GROUP" or "APPLICATION". This allows for metric per server, per server group or on the whole application.
Optional:
	--period
		The amount of time (in seconds) between 2 data points. [default: 60]
	--description
		Description for the metric. [default: ""]
	--unit
		The unit for the metric. This is shown on the axis in the widgets. [default: ""]
	--attachTo
		Describes what the relation of this Metric is. Options are SERVER, SERVERGROUP, APPLICATION, REQUEST, DATABASE, QUERY and ANALYSIS.
	--source
		Describes who added the metric. Can be chosen by the user. [default: "cli"]
`,
		Run: func(cmd *Command, args []string) {
			var name, description, dataType, subject, unit, attachTo, source string
			var period int
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the metric.")
			cmd.Flag.StringVar(&description, "description", "", "Description for the metric.")
			cmd.Flag.StringVar(&dataType, "dataType", DEFAULT_STRING_FLAG_VALUE, `The following data types are defined: "LONG", "DOUBLE", "HISTOGRAM".`)
			cmd.Flag.StringVar(&subject, "subject", DEFAULT_STRING_FLAG_VALUE, `A metric is defined on either a "SERVER", "GROUP" or "APPLICATION".`)
			cmd.Flag.StringVar(&unit, "unit", "", "The unit for the metric.")
			cmd.Flag.StringVar(&attachTo, "attachTo", "", "Describes what the relation of this Metric is.")
			cmd.Flag.StringVar(&source, "source", "cli", "Describes who added the metric.")
			cmd.Flag.IntVar(&period, "period", 60, "The amount of time (in seconds) between 2 data points.")
			cmd.ParseArgs(args)

			if name == DEFAULT_STRING_FLAG_VALUE || dataType == DEFAULT_STRING_FLAG_VALUE || subject == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			cmd.PrintResult(cmd.Capi.CreateMetric(name, description, dataType, unit, subject, source, period))
		},
	},
	{
		Name:      "update",
		UsageLine: `metric update (--name | --id) [--description --dataType --subject --unit --period --attachTo --source]`,
		Long: `
Update a CoScale metric object.

The flags for update metric action are:

Mandatory:
	--name 
		specify the name of the metric.
	--id
		Unique identifier, if we want to update the name of a metric, this become mandatory
Optional:
	--description
			Description for the metric.
	--dataType
			The following data types are defined: "LONG", "DOUBLE", "HISTOGRAM".
	--subject
			A metric is defined on either a "SERVER", "GROUP" or "APPLICATION". This allows for metric per server, per server group or on the whole application.
	--unit
			The unit for the metric. This is shown on the axis in the widgets. [default: ""]
	--attachTo
			Describes what the relation of this Metric is. Options are SERVER, SERVERGROUP, APPLICATION, REQUEST, DATABASE, QUERY and ANALYSIS.
	--source
			Describes who added the metric. Can be chosen by the user. [default: "cli"]
	--period
			The amount of time (in seconds) between 2 data points. [default: 60]
`,
		Run: func(cmd *Command, args []string) {
			var name, description, dataType, subject, unit, attachTo, source string
			var period int
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier")
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the metric.")
			cmd.Flag.StringVar(&description, "description", DEFAULT_STRING_FLAG_VALUE, "Description for the metric.")
			cmd.Flag.StringVar(&dataType, "dataType", DEFAULT_STRING_FLAG_VALUE, `The following data types are defined: "LONG", "DOUBLE", "HISTOGRAM".`)
			cmd.Flag.StringVar(&subject, "subject", "", `A metric is defined on either a "SERVER", "GROUP" or "APPLICATION".`)
			cmd.Flag.StringVar(&unit, "unit", DEFAULT_STRING_FLAG_VALUE, "The unit for the metric.")
			cmd.Flag.StringVar(&attachTo, "attachTo", "", "Describes what the relation of this Metric is.")
			cmd.Flag.StringVar(&source, "source", DEFAULT_STRING_FLAG_VALUE, "Describes who added the metric.")
			cmd.Flag.IntVar(&period, "period", -1, "The amount of time (in seconds) between 2 data points.")
			cmd.ParseArgs(args)

			var metricObj = &api.Metric{}
			var err error
			if name == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if id != -1 {
				err = cmd.Capi.GetObjectRef("metric", id, metricObj)
			} else {
				err = cmd.Capi.GetObejctRefByName("metric", name, metricObj)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}
			//update the metric object values
			if name != DEFAULT_STRING_FLAG_VALUE {
				metricObj.Name = name
			}
			if dataType != DEFAULT_STRING_FLAG_VALUE {
				metricObj.DataType = dataType
			}
			if description != DEFAULT_STRING_FLAG_VALUE {
				metricObj.Description = description
			}
			if period != -1 {
				metricObj.Period = period
			}
			if source != DEFAULT_STRING_FLAG_VALUE {
				metricObj.Source = source
			}
			if unit != DEFAULT_STRING_FLAG_VALUE {
				metricObj.Unit = unit
			}

			cmd.PrintResult(cmd.Capi.UpdateMetric(metricObj))
		},
	},
}

var MetricGroupObject = NewCommand("metricgroup", "metricgroup <action> [--<field>='<data>']", MetricGroupActions)

var MetricGroupActions = []*Command{
	ListCmd("metricgroup"),
	GetCmd("metricgroup"),
	AddObjToGroupCmd("metric", &api.Metric{}, &api.MetricGroup{}),
	DeleteObjFromGroupCmd("metric", &api.Metric{}, &api.MetricGroup{}),
	{
		Name:      "new",
		UsageLine: `servergroup new (--name --subject) [--description --state --source]`,
		Long: `
Create a new CoScale metricgroup object.

The flags for new metricgroup action are:

Mandatory:
	--name 
		Name for the metric group.
	--subject
		The subject type of the metric group. "APPLICATION", "SERVERGROUP" or "SERVER".
Optional:
	--description
		Description for the metric group.
	--type
		Describes the type of metric group.
	--state
		"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.
	--source
		Describes who added the metric group. Can be chosen by the user. [default: "cli"]
`,
		Run: func(cmd *Command, args []string) {
			var name, description, Type, state, source, subject string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the metric group.")
			cmd.Flag.StringVar(&description, "description", "", "Description for the metric group.")
			cmd.Flag.StringVar(&Type, "type", "", "Describes the type of metric group.")
			cmd.Flag.StringVar(&subject, "subject", DEFAULT_STRING_FLAG_VALUE, `The subject type of the metric group. "APPLICATION", "SERVERGROUP" or "SERVER".`)
			cmd.Flag.StringVar(&state, "state", "ENABLED", `"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.`)
			cmd.Flag.StringVar(&source, "source", "cli", "Describes who added the metric group.")
			cmd.ParseArgs(args)

			if name == DEFAULT_STRING_FLAG_VALUE || subject == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			cmd.PrintResult(cmd.Capi.CreateMetricGroup(name, description, Type, state, subject, source))
		},
	},
	{
		Name:      "update",
		UsageLine: `metricgroup update (--name | --id) [--description --type --state --source]`,
		Long: `
Update a CoScale metricgroup object.

The flags for update metricgroup action are:

Mandatory:
	--name 
		Specify the name of the metricgroup.
Optional:
	--id
		Unique identifier, if we want to update the name of the metricgroup, this become mandatory.
	--description
		Description for the metric group.
	--type
		Describes the type of metric group.
	--state
		"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.
	--source
		Describes who added the metric group. Can be chosen by the user. [default: "cli"]
`,
		Run: func(cmd *Command, args []string) {
			var id int64
			var name, description, Type, state, source string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier")
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the metric group.")
			cmd.Flag.StringVar(&description, "description", DEFAULT_STRING_FLAG_VALUE, "Description for the metric group.")
			cmd.Flag.StringVar(&Type, "type", DEFAULT_STRING_FLAG_VALUE, "Describes the type of metric group.")
			cmd.Flag.StringVar(&state, "state", DEFAULT_STRING_FLAG_VALUE, `"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.`)
			cmd.Flag.StringVar(&source, "source", DEFAULT_STRING_FLAG_VALUE, "Describes who added the metric group.")
			cmd.ParseArgs(args)

			if name == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}

			var metricGroupObj = &api.MetricGroup{}
			var err error
			if name == "" {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if id != -1 {
				err = cmd.Capi.GetObjectRef("metricgroup", id, metricGroupObj)
			} else {
				err = cmd.Capi.GetObejctRefByName("metricgroup", name, metricGroupObj)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}

			//update the metricgroup object values
			if name != DEFAULT_STRING_FLAG_VALUE {
				metricGroupObj.Name = name
			}
			if description != DEFAULT_STRING_FLAG_VALUE {
				metricGroupObj.Description = description
			}
			if Type != DEFAULT_STRING_FLAG_VALUE {
				metricGroupObj.Type = Type
			}
			if source != DEFAULT_STRING_FLAG_VALUE {
				metricGroupObj.Source = source
			}
			if state != DEFAULT_STRING_FLAG_VALUE {
				metricGroupObj.State = state
			}
			cmd.PrintResult(cmd.Capi.UpdateMetricGroup(metricGroupObj))
		},
	},
}
