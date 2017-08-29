package command

import (
	"coscale/api"
	"encoding/json"
	"os"
)

// metricSubCommands will contain subcommands for metric command and also actions for it.
var metricSubCommands = append(MetricActions, []*Command{
	// subcommands of metric
	dimensionObject,
}...)

/**
 * Metric Actions
 */

// MetricObject defines the metric command on the CLI.
var MetricObject = NewCommand("metric", "metric <action> [--<field>='<data>']", metricSubCommands)

// MetricActions defines the metric actions on the CLI.
var MetricActions = []*Command{
	ListCmd("metric"),
	GetCmd("metric"),
	DeleteCmd(&api.Metric{}, "metric"),
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
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier.")
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
		UsageLine: `metric new (--name --dataType --subject) [--period --description --unit --attachTo]`,
		Long: `
Create a new CoScale metric object.

The flags for new metric action are:

Mandatory:
	--name
		specify the name of the metric.
	--dataType
		The following data types are defined: "LONG", "DOUBLE", "HISTOGRAM", "COUNT", "COUNTER", "BINARY".
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
			cmd.Flag.StringVar(&source, "source", "cli", "Deprecated.")
			cmd.Flag.IntVar(&period, "period", 60, "The amount of time (in seconds) between 2 data points.")
			cmd.ParseArgs(args)

			if name == DEFAULT_STRING_FLAG_VALUE || dataType == DEFAULT_STRING_FLAG_VALUE || subject == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}

			cmd.PrintResult(cmd.Capi.CreateMetric(name, description, dataType, unit, subject, period))
		},
	},
	{
		Name:      "update",
		UsageLine: `metric update (--name | --id) [--description --dataType --subject --unit --period --attachTo]`,
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
	--period
			The amount of time (in seconds) between 2 data points. [default: 60]
`,
		Run: func(cmd *Command, args []string) {
			var name, description, dataType, subject, unit, attachTo, source string
			var period int
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier.")
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the metric.")
			cmd.Flag.StringVar(&description, "description", DEFAULT_STRING_FLAG_VALUE, "Description for the metric.")
			cmd.Flag.StringVar(&dataType, "dataType", DEFAULT_STRING_FLAG_VALUE, `The following data types are defined: "LONG", "DOUBLE", "HISTOGRAM".`)
			cmd.Flag.StringVar(&subject, "subject", DEFAULT_STRING_FLAG_VALUE, `A metric is defined on either a "SERVER", "GROUP" or "APPLICATION".`)
			cmd.Flag.StringVar(&unit, "unit", DEFAULT_STRING_FLAG_VALUE, "The unit for the metric.")
			cmd.Flag.StringVar(&attachTo, "attachTo", "", "Describes what the relation of this Metric is.")
			cmd.Flag.StringVar(&source, "source", DEFAULT_STRING_FLAG_VALUE, "Deprecated.")
			cmd.Flag.IntVar(&period, "period", -1, "The amount of time (in seconds) between 2 data points.")
			cmd.ParseArgs(args)

			var err error
			var metricObj = &api.Metric{}
			if id != -1 {
				err = cmd.Capi.GetObjectRef("metric", id, metricObj)
			} else if name != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName("metric", name, metricObj)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
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
			if subject != DEFAULT_STRING_FLAG_VALUE {
				metricObj.Subject = subject
			}
			if unit != DEFAULT_STRING_FLAG_VALUE {
				metricObj.Unit = unit
			}

			cmd.PrintResult(cmd.Capi.UpdateMetric(metricObj))
		},
	},
}

// MetricGroupObject defines the metric group command on the CLI.
var MetricGroupObject = NewCommand("metricgroup", "metricgroup <action> [--<field>='<data>']", MetricGroupActions)

// MetricGroupActions defines the metric group actions on the CLI.
var MetricGroupActions = []*Command{
	ListCmd("metricgroup"),
	GetCmd("metricgroup"),
	DeleteCmd(&api.MetricGroup{}, "metricgroup"),
	AddObjToGroupCmd("metric", &api.Metric{}, &api.MetricGroup{}),
	DeleteObjFromGroupCmd("metric", &api.Metric{}, &api.MetricGroup{}),
	{
		Name:      "new",
		UsageLine: `metricgroup new (--name --subject) [--description --state]`,
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
`,
		Run: func(cmd *Command, args []string) {
			var name, description, Type, state, source, subject string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the metric group.")
			cmd.Flag.StringVar(&description, "description", "", "Description for the metric group.")
			cmd.Flag.StringVar(&Type, "type", "", "Describes the type of metric group.")
			cmd.Flag.StringVar(&subject, "subject", DEFAULT_STRING_FLAG_VALUE, `The subject type of the metric group. "APPLICATION", "SERVERGROUP" or "SERVER".`)
			cmd.Flag.StringVar(&state, "state", "ENABLED", `"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.`)
			cmd.Flag.StringVar(&source, "source", "cli", "Deprecated.")
			cmd.ParseArgs(args)

			if name == DEFAULT_STRING_FLAG_VALUE || subject == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			cmd.PrintResult(cmd.Capi.CreateMetricGroup(name, description, Type, state, subject))
		},
	},
	{
		Name:      "update",
		UsageLine: `metricgroup update (--name | --id) [--description --type --state]`,
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
`,
		Run: func(cmd *Command, args []string) {
			var id int64
			var name, description, Type, state, source string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier.")
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the metric group.")
			cmd.Flag.StringVar(&description, "description", DEFAULT_STRING_FLAG_VALUE, "Description for the metric group.")
			cmd.Flag.StringVar(&Type, "type", DEFAULT_STRING_FLAG_VALUE, "Describes the type of metric group.")
			cmd.Flag.StringVar(&state, "state", DEFAULT_STRING_FLAG_VALUE, `"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.`)
			cmd.Flag.StringVar(&source, "source", DEFAULT_STRING_FLAG_VALUE, "Deprecated.")
			cmd.ParseArgs(args)

			var err error
			var metricGroupObj = &api.MetricGroup{}
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
			if state != DEFAULT_STRING_FLAG_VALUE {
				metricGroupObj.State = state
			}
			cmd.PrintResult(cmd.Capi.UpdateMetricGroup(metricGroupObj))
		},
	},
}

/**
 * Dimension subcommand with actions.
 */
// Type will be a new subcommand of alert and it will have actions.
var dimensionObjectName = "dimension"
var dimensionObject = NewCommand(dimensionObjectName, "dimension <action> [--<field>='<data>']", dimensionActions)

// DimensionActions will contain the actions for metric dimension subcommand.
var dimensionActions = []*Command{
	{
		Name:      "new",
		UsageLine: `metric dimension new (--name) [--id|--metric]`,
		Long: `
Create a new CoScale dimension object for a metric.

Metric dimensions enables us to show metrics at different levels. For example for RabbitMQ
we want to show the total number of queued messages, but we also want to be able to split these into the number of queued messages per queue.

The flags for dimension new action are:

Mandatory:
	--name
		Specify the name of the new metric dimension.
Optional:
	--id
		Unique identifier for the metric.
	or
	--metric
		Specify the name of the metric.
`,
		Run: func(cmd *Command, args []string) {
			var name, metric string
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Specify the name of the new metric dimension.")
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier for the metric.")
			cmd.Flag.StringVar(&metric, "metric", DEFAULT_STRING_FLAG_VALUE, "Specify the name of the metric.")

			cmd.ParseArgs(args)

			if name == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}

			var dimension string
			var err error
			// Create dimension.
			dimension, err = cmd.Capi.CreateDimension(name)
			if err != nil {
				cmd.PrintResult("", err)
			}
			var dimensionObj *api.Dimension
			if err := json.Unmarshal([]byte(dimension), &dimensionObj); err != nil {
				cmd.PrintResult("", err)
			}

			// Get the metric.
			if id == -1 && metric != DEFAULT_STRING_FLAG_VALUE {
				var metricObj = &api.Metric{}
				err = cmd.Capi.GetObejctRefByName("metric", metric, metricObj)
				if err != nil {
					cmd.PrintResult("", err)
				}

				id = metricObj.ID
			}

			// if no metric to asociate with the dimension do not continue.
			if id == -1 {
				cmd.PrintResult(dimension, nil)
			}
			// Associate dimension with the metric.
			cmd.PrintResult(cmd.Capi.AddMetricDimension(id, dimensionObj.ID))
		},
	},
	{
		Name:      "list",
		UsageLine: `metric dimension list (--metric|--metricId)`,
		Long: `
Get all the dimensions for a metric

The flags for dimension list action are:

Mandatory:
	--metric
		Specify the name of the metric.
or
	--metricId
		Unique identifier for the metric.
`,
		Run: func(cmd *Command, args []string) {
			var metric string
			var metricID int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&metric, "metric", DEFAULT_STRING_FLAG_VALUE, "Specify the name of the metric.")
			cmd.Flag.Int64Var(&metricID, "metricId", -1, "Unique identifier for the metric.")

			cmd.ParseArgs(args)

			var metricObj = &api.Metric{}
			var err error
			if metricID != -1 {
				err = cmd.Capi.GetObjectRef("metric", metricID, metricObj)
			} else if metric != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName("metric", metric, metricObj)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}

			cmd.PrintResult(cmd.Capi.GetDimensions(metricObj.ID))
		},
	},
}
