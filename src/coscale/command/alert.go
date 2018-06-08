package command

import (
	"coscale/api"
	"os"
)

// alertSubCommands will contain subcommands for alert command and also actions for it.
var alertSubCommands = append(AlertActions, []*Command{
	// subcommands of alert
	alertTypeObject,
	alertTriggerObject,
}...)

/**
 *
 * Alert actions.
 *
 */

var alertObjectName = "alert"

// AlertObject defines the alert command on the CLI.
var AlertObject = NewCommand(alertObjectName, "alert <action> [--<field>='<data>']", alertSubCommands)

// AlertActions defines the alert actions on the CLI.
var AlertActions = []*Command{
	{
		Name:      "list",
		UsageLine: "alert list [--filter]",
		Long: `
Get all alerts from CoScale Api.

The flags for list alert action are:

Optional:
	--filter
		List actions filtered by unresolved or by unacknowledged.
`,
		Run: func(cmd *Command, args []string) {
			var filter string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&filter, "filter", DEFAULT_STRING_FLAG_VALUE, "List actions filtered by unresolved or by unacknowledged.")
			cmd.ParseArgs(args)

			switch filter {
			case DEFAULT_STRING_FLAG_VALUE:
				cmd.PrintResult(cmd.Capi.GetObjects(alertObjectName))
			case "unresolved":
				cmd.PrintResult(cmd.Capi.GetAlertsBy("selectByResolved"))
			case "unacknowledged":
				cmd.PrintResult(cmd.Capi.GetAlertsBy("selectByAcknowledged"))
			default:
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
		},
	},
	{
		Name:      "acknowledge",
		UsageLine: "alert acknowledge (--id)",
		Long: `
Acknowledge an alert.

The flags for acknowledge alert action are:
Mandatory:
	--id
		The id of the alert.
`,
		Run: func(cmd *Command, args []string) {
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&id, "id", -1, "The id of the alert.")
			cmd.ParseArgs(args)
			var alert = &api.Alert{}
			var err error
			if id == -1 {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err = cmd.Capi.GetObjectRef("alert", id, alert); err != nil {
				cmd.PrintResult("", err)
			}
			cmd.PrintResult(cmd.Capi.AlertSolution(alert, "acknowledge"))
		},
	},
	{
		Name:      "resolve",
		UsageLine: "alert resolve (--id)",
		Long: `
Resolve an alert.

The flags for resolve alert action are:
Mandatory:
	--id
		The id of the alert.
`,
		Run: func(cmd *Command, args []string) {
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&id, "id", -1, "The id of the alert.")
			cmd.ParseArgs(args)
			var alert = &api.Alert{}
			var err error
			if id == -1 {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err = cmd.Capi.GetObjectRef("alert", id, alert); err != nil {
				cmd.PrintResult("", err)
			}
			cmd.PrintResult(cmd.Capi.AlertSolution(alert, "resolve"))
		},
	},
}

/**
 *
 * AlertType subcommand with actions.
 *
 */

// Type will be a new subcommand of alert and it will have actions.
var alertTypeObjectName = "type"
var alertTypeObject = NewCommand(alertTypeObjectName, "type <action> [--<field>='<data>']", alertTypeActions)

// AlertTypeActions will contain the actions for alert type subcommand.
var alertTypeActions = []*Command{
	GetCmd("alerttype", "type"),
	ListCmd("alerttype", "type"),
	{
		Name:      "new",
		UsageLine: `alert type new (--name --handle) [--description --backupHandle --backupSeconds --escalationHandle --escalationSeconds]`,
		Long: `
Create a new CoScale alert type.

The flags for new type action are:

Mandatory:
	--name
		Name for the new alert type.
	--handle
		The handle fields describe how an alert is delivered to the user.
		Is a list of objects, each object describes a delivery mechanism.
		At the moment we support sending an email to a user, sending an email to an email address or integrations
		for third party services:
		e.g.
		--handle "EMAIL:support@coscale.com"
		also multiple contacts can be provided
		--handle "EMAILUSER:1 EMAIL:support@coscale.com SLACK:https://hooks.slack.com..."
Optional:
	--description
		Description for the alert type.
	--backupHandle
		AlertType can have 3 levels of handlers set. First an alert is sent. If there is no response within backupSeconds,
		a backup-alert is sent. If there is no response within escalationSeconds, an escalation is sent.
	--backupSeconds
		Number of second to wait until notifications are sent to the second handle level.
	--escalationHandle
		Third handle level.
	--escalationSeconds
		Number of second to wait until notifications are sent to the third handle level.
`,
		Run: func(cmd *Command, args []string) {
			var name, handle, description, backupHandle, escalationHandle, source string
			var backupSeconds, escalationSeconds int64

			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the new alert type.")
			cmd.Flag.StringVar(&handle, "handle", DEFAULT_STRING_FLAG_VALUE, "The handle fields describe how an alert is delivered to the user.")
			cmd.Flag.StringVar(&description, "description", "", "Description for the alert type.")
			cmd.Flag.StringVar(&backupHandle, "backupHandle", DEFAULT_STRING_FLAG_VALUE, "The handle fields describe how an alert is delivered to the user.")
			cmd.Flag.Int64Var(&backupSeconds, "backupSeconds", -1, "Number of second to wait until notifications are sent to the second handle level.")
			cmd.Flag.StringVar(&escalationHandle, "escalationHandle", DEFAULT_STRING_FLAG_VALUE, "The handle fields describe how an alert is delivered to the user.")
			cmd.Flag.Int64Var(&escalationSeconds, "escalationSeconds", -1, "Number of second to wait until notifications are sent to the third handle level.")
			cmd.Flag.StringVar(&source, "source", "cli", "Deprecated.")
			cmd.ParseArgs(args)

			// Check if values were provided for mandatory flags.
			if name == DEFAULT_STRING_FLAG_VALUE || handle == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}

			// Parse the alert handle.
			var err error
			handle, err = api.ParseHandle(handle)
			if err != nil {
				cmd.PrintResult("", err)
			}

			// Parse the alert backupHandle.
			if backupHandle != DEFAULT_STRING_FLAG_VALUE {
				backupHandle, err = api.ParseHandle(backupHandle)
				if err != nil {
					cmd.PrintResult("", err)
				}
			}

			// Parse the alert escalationHandle.
			if escalationHandle != DEFAULT_STRING_FLAG_VALUE {
				escalationHandle, err = api.ParseHandle(escalationHandle)
				if err != nil {
					cmd.PrintResult("", err)
				}
			}

			cmd.PrintResult(cmd.Capi.CreateType(name, description, handle, backupHandle, escalationHandle, backupSeconds, escalationSeconds))
		},
	},
	{
		Name:      "update",
		UsageLine: `alert type update (--name | --id) [--name --handle --description --backupHandle --backupSeconds --escalationHandle --escalationSeconds]`,
		Long: `
Update an existing CoScale alert type.

The flags for update type action are:

Mandatory:
	--name
		Name for the alert type.
Optional:
	--id
		Unique identifier, if we want to update the name of the alert type, this become mandatory.
	--handle
		The handle fields describe how an alert is delivered to the user.
		Is a list of objects, each object describes a delivery mechanism.
		At the moment we support sending an email to a user, sending an email to an email address or integrations
		for third party services:
		e.g.
		--handle "EMAIL:support@coscale.com"
		also multiple contacts can be provided
		--handle "EMAILUSER:1 EMAIL:support@coscale.com SLACK:https://hooks.slack.com..."
	--description
		Description for the alert type.
	--backupHandle
		AlertType can have 3 levels of handlers set. First an alert is sent. If there is no response within backupSeconds,
		a backup-alert is sent. If there is no response within escalationSeconds, an escalation is sent.
	--backupSeconds
		Number of second to wait until notifications are sent to the second handle level.
	--escalationHandle
		Third handle level.
	--escalationSeconds
		Number of second to wait until notifications are sent to the third handle level.
`,
		Run: func(cmd *Command, args []string) {
			var name, handle, description, backupHandle, escalationHandle, source string
			var id, backupSeconds, escalationSeconds int64

			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the new alert type.")
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier.")
			cmd.Flag.StringVar(&handle, "handle", DEFAULT_STRING_FLAG_VALUE, "The handle fields describe how an alert is delivered to the user.")
			cmd.Flag.StringVar(&description, "description", DEFAULT_STRING_FLAG_VALUE, "Description for the alert type.")
			cmd.Flag.StringVar(&backupHandle, "backupHandle", DEFAULT_STRING_FLAG_VALUE, "The handle fields describe how an alert is delivered to the user.")
			cmd.Flag.Int64Var(&backupSeconds, "backupSeconds", -1, "Number of second to wait until notifications are sent to the second handle level.")
			cmd.Flag.StringVar(&escalationHandle, "escalationHandle", DEFAULT_STRING_FLAG_VALUE, "The handle fields describe how an alert is delivered to the user.")
			cmd.Flag.Int64Var(&escalationSeconds, "escalationSeconds", -1, "Number of second to wait until notifications are sent to the third handle level.")
			cmd.Flag.StringVar(&source, "source", DEFAULT_STRING_FLAG_VALUE, "Deprecated.")
			cmd.ParseArgs(args)

			var err error
			// Get the existing trigger and create the update object.
			var alertTypeObj = &api.AlertType{}
			if id != -1 {
				err = cmd.Capi.GetObjectRef("alerttype", id, alertTypeObj)
			} else if name != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObjectRefByName("alerttype", name, alertTypeObj)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}

			if err != nil {
				cmd.PrintResult("", err)
			}

			// update the alertType object values
			if name != DEFAULT_STRING_FLAG_VALUE {
				alertTypeObj.Name = name
			}
			if description != DEFAULT_STRING_FLAG_VALUE {
				alertTypeObj.Description = description
			}
			if handle != DEFAULT_STRING_FLAG_VALUE {
				handle, err = api.ParseHandle(handle)
				if err != nil {
					cmd.PrintResult("", err)
				}
				alertTypeObj.Handle = handle
			}
			if backupHandle != DEFAULT_STRING_FLAG_VALUE {
				backupHandle, err = api.ParseHandle(backupHandle)
				if err != nil {
					cmd.PrintResult("", err)
				}
				alertTypeObj.BackupHandle = backupHandle
			}
			if backupSeconds != -1 {
				alertTypeObj.BackupSeconds = backupSeconds
			}
			if escalationHandle != DEFAULT_STRING_FLAG_VALUE {
				escalationHandle, err = api.ParseHandle(escalationHandle)
				if err != nil {
					cmd.PrintResult("", err)
				}
				alertTypeObj.EscalationHandle = escalationHandle
			}
			if escalationSeconds != -1 {
				alertTypeObj.EscalationSeconds = escalationSeconds
			}

			cmd.PrintResult(cmd.Capi.UpdateType(alertTypeObj))
		},
	},
	DeleteCmd(&api.AlertType{}, "alerttype", "type"),
}

/**
 *
 * Trigger subcommand with actions.
 *
 */

// Type will be a new subcommand of alert and it will have actions.
var alertTriggerObjectName = "trigger"
var alertTriggerObject = NewCommand(alertTriggerObjectName, "trigger <action> [--<field>='<data>']", alertTriggerActions)

// AlertTypeActions will contain the actions for alert type subcommand.
var alertTriggerActions = []*Command{
	{
		Name:      "list",
		UsageLine: "alert trigger list (--id | --name)",
		Long: `
Get all alert triggers for an alert type from CoScale Api.

The flags for list trigger action are:

Mandatory:
	--name
		specify the name of the alert type for triggers.
	or
	--id
		specify the alert type id for triggers.
`,
		Run: func(cmd *Command, args []string) {
			var id int64
			var name string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier of alert type.")
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name of the alert type.")

			cmd.ParseArgs(args)

			var alertTypeObj = &api.AlertType{}
			var err error
			if id != -1 {
				cmd.PrintResult(cmd.Capi.GetTriggers(id))
			} else if name != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObjectRefByName("alerttype", name, alertTypeObj)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}

			cmd.PrintResult(cmd.Capi.GetTriggers(alertTypeObj.ID))
		},
	},
	{
		Name:      "new",
		UsageLine: `alert trigger new (--name --config --metric|--metricid) [--autoresolve --typename|--typeid --description --server|--serverid --servergroup|--servergroupid]`,
		Long: `
Create a new CoScale alert trigger.

The flags for new trigger action are:

Mandatory:
	--name
		Name for the new trigger.
	--config
		The trigger configuration which is formatted as follows:
			For metrics with DataType DOUBLE:
				avg(300) > 25 (if the average value over 5 minutes is larger than 25, trigger an alert.)
			For metrics with DataType HISTOGRAM:
				avg(99, 300) >= 50 (if the average of the 99th percentile over 5 minutes is larger or equal to 50, trigger an alert.)
	--metric
		The name of the metric which will be the subject of the alert.
	or
	--metricid
		The id of the metric which will be the subject of the alert.
Optional:
	--autoresolve
		The amount of seconds to wait until the alert will be auto-resolved [default: null]
	--typename
		specify the name of the alert type for triggers. [default: "Default alerts"]
	or
	--typeid
		specify the alert type id for triggers. [default: "Default alerts"]
	--description
		Description for the alert trigger.
	--server
		The server name for which the alert will be triggered.
	or
	--serverid
		The server id for which the alert will be triggered.
	--servergroup
		The servergroup name for which the alert will be triggered.
	or
	--servergroupid
		The servergroup id for which the alert will be triggered.
	Note: if no server or servergroup is provided the trigger will be set for the entire application.
`,
		Run: func(cmd *Command, args []string) {
			var name, config, metric, description, server, serverGroup, source, typeName string
			var metricID, autoResolve, serverID, serverGroupID, typeID int64
			var onApp bool

			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the new trigger.")
			cmd.Flag.StringVar(&config, "config", DEFAULT_STRING_FLAG_VALUE, "The trigger configuration.")
			cmd.Flag.Int64Var(&autoResolve, "autoresolve", -1, "The amount of seconds to wait until the alert will be auto-resolved.")
			cmd.Flag.StringVar(&metric, "metric", DEFAULT_STRING_FLAG_VALUE, "The name of the metric which will be the subject of the alert.")
			cmd.Flag.Int64Var(&metricID, "metricid", -1, "The id of the metric which will be the subject of the alert.")
			cmd.Flag.StringVar(&description, "description", "", "Description for the alert trigger.")
			cmd.Flag.StringVar(&server, "server", DEFAULT_STRING_FLAG_VALUE, "The server name for which the alert will be triggered.")
			cmd.Flag.Int64Var(&serverID, "serverid", -1, "The server id for which the alert will be triggered.")
			cmd.Flag.StringVar(&serverGroup, "servergroup", DEFAULT_STRING_FLAG_VALUE, "The servergroup name for which the alert will be triggered.")
			cmd.Flag.Int64Var(&serverGroupID, "servergroupid", -1, "The server id for which the alert will be triggered.")
			cmd.Flag.StringVar(&source, "source", "cli", "Deprecated.")
			cmd.Flag.StringVar(&typeName, "typename", "Default alerts", "Specify the name of the alert type for triggers.")
			cmd.Flag.Int64Var(&typeID, "typeid", -1, "Specify the alert type id for triggers.")

			cmd.ParseArgs(args)

			// Check if values were provided for mandatory flags.
			if name == DEFAULT_STRING_FLAG_VALUE || config == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if metric == DEFAULT_STRING_FLAG_VALUE && metricID == -1 {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}

			// Get the metric id
			var metricObj = &api.Metric{}
			var err error
			if metricID == -1 {
				err = cmd.Capi.GetObjectRefByName("metric", metric, metricObj)
				if err != nil {
					cmd.PrintResult("", err)
				}
				// if didn't exit due to error...
				metricID = metricObj.ID
			}

			// Get the server id
			var serverObj = &api.Server{}
			if serverID == -1 && server != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObjectRefByName("server", server, serverObj)
				if err != nil {
					cmd.PrintResult("", err)
				}
				// if didn't exit due to error...
				serverID = serverObj.ID
			}

			// Get the servergroup id
			var serverGroupObj = &api.ServerGroup{}
			if serverGroupID == -1 && serverGroup != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObjectRefByName("servergroup", serverGroup, serverGroupObj)
				if err != nil {
					cmd.PrintResult("", err)
				}
				// if didn't exit due to error...
				serverGroupID = serverGroupObj.ID
			}

			// if no error and no server or servergroup id was found then the trigger is for app.
			onApp = serverID == -1 && serverGroupID == -1

			// get the alert type for the trigger
			var alertTypeObj = &api.AlertType{}
			if typeID == -1 {
				err = cmd.Capi.GetObjectRefByName("alerttype", typeName, alertTypeObj)
				if err != nil {
					cmd.PrintResult("", err)
				}
				// if didn't exit due to error...
				typeID = alertTypeObj.ID
			}

			cmd.PrintResult(cmd.Capi.CreateTrigger(name, description, config, typeID, autoResolve, metricID, serverID, serverGroupID, onApp))
		},
	},
	{
		Name:      "update",
		UsageLine: `alert trigger update (--typeid --id|--typename --name) [--autoresolve --name --config --metric|--metricid --description --server|--serverid --servergroup|--servergroupid]`,
		Long: `
Update a existing CoScale alert trigger.

The flags for update trigger action are:

Mandatory
	--typeid
		Specify the alert type id for the trigger.
	--id
		Unique identifier, if we want to update the name of the trigger, this become mandatory.
	or
	--typename
		Specify the name of the alert type for the trigger.
	--name
		Name for the trigger.
Optional:
	--autoresolve
		The amount of seconds to wait until the alert will be auto-resolved. [default: null]
	--config
		The trigger configuration which is formatted as follows:
			For metrics with DataType DOUBLE:
				avg(300) > 25 (if the average value over 5 minutes is larger than 25, trigger an alert.)
			For metrics with DataType HISTOGRAM:
				avg(99, 300) >= 50 (if the average of the 99th percentile over 5 minutes is larger or equal to 50, trigger an alert.)
	--metric
		The name of the metric which will be the subject of the alert.
	or
	--metricid
		The id of the metric which will be the subject of the alert.
	--description
		Description for the alert trigger.
	--server
		The server name for which the alert will be triggered.
	or
	--serverid
		The server id for which the alert will be triggered.
	--servergroup
		The servergroup name for which the alert will be triggered.
	or
	--servergroupid
		The servergroup id for which the alert will be triggered.
	Note: if no server or servergroup is provided the tigger will be set for entire application.
`,
		Run: func(cmd *Command, args []string) {
			var name, config, metric, description, server, serverGroup, source, typeName string
			var id, autoResolve, metricID, serverID, serverGroupID, typeID int64
			var onApp bool

			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the new trigger.")
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier for trigger.")
			cmd.Flag.StringVar(&config, "config", DEFAULT_STRING_FLAG_VALUE, "The trigger configuration.")
			cmd.Flag.Int64Var(&autoResolve, "autoresolve", -1, "The amount of seconds to wait until the alert will be auto-resolved.")
			cmd.Flag.StringVar(&metric, "metric", DEFAULT_STRING_FLAG_VALUE, "The name of the metric which will be the subject of the alert.")
			cmd.Flag.Int64Var(&metricID, "metricid", -1, "The id of the metric which will be the subject of the alert.")
			cmd.Flag.StringVar(&description, "description", DEFAULT_STRING_FLAG_VALUE, "Description for the alert trigger.")
			cmd.Flag.StringVar(&server, "server", DEFAULT_STRING_FLAG_VALUE, "The server name for which the alert will be triggered.")
			cmd.Flag.Int64Var(&serverID, "serverid", -1, "The server id for which the alert will be triggered.")
			cmd.Flag.StringVar(&serverGroup, "servergroup", DEFAULT_STRING_FLAG_VALUE, "The servergroup name for which the alert will be triggered.")
			cmd.Flag.Int64Var(&serverGroupID, "servergroupid", -1, "The server id for which the alert will be triggered.")
			cmd.Flag.StringVar(&source, "source", DEFAULT_STRING_FLAG_VALUE, "Deprecated.")
			cmd.Flag.StringVar(&typeName, "typename", DEFAULT_STRING_FLAG_VALUE, "Specify the name of the alert type for triggers.")
			cmd.Flag.Int64Var(&typeID, "typeid", -1, "Specify the alert type id for triggers.")

			cmd.ParseArgs(args)

			// Check if values were provided for mandatory flags.
			if id == -1 && name == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}

			if typeID == -1 && typeName == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}

			// Get the metric id
			var metricObj = &api.Metric{}
			var err error
			if metricID == -1 && metric != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObjectRefByName("metric", metric, metricObj)
				if err != nil {
					cmd.PrintResult("", err)
				}
				// if didn't exit due to error...
				metricID = metricObj.ID
			}

			// Get the server id
			var serverObj = &api.Server{}
			if serverID == -1 && server != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObjectRefByName("server", server, serverObj)
				if err != nil {
					cmd.PrintResult("", err)
				}
				// if didn't exit due to error...
				serverID = serverObj.ID
			}

			// Get the servergroup id
			var serverGroupObj = &api.ServerGroup{}
			if serverGroupID == -1 && serverGroup != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObjectRefByName("servergroup", serverGroup, serverGroupObj)
				if err != nil {
					cmd.PrintResult("", err)
				}
				// if didn't exit due to error...
				serverGroupID = serverGroupObj.ID
			}

			// get the alert type for the trigger
			var alertTypeObj = &api.AlertType{}
			if typeID == -1 {
				err = cmd.Capi.GetObjectRefByName("alerttype", typeName, alertTypeObj)
				if err != nil {
					cmd.PrintUsage()
					os.Exit(EXIT_FLAG_ERROR)
				}
				// if didn't exit due to error...
				typeID = alertTypeObj.ID
			}

			// Get the existing trigger and create the update object.
			var alertTriggerObj = &api.AlertTrigger{}
			if id != -1 {
				err = cmd.Capi.GetObjectRefFromGroup("alerttype", "trigger", typeID, id, alertTriggerObj)
			} else if name != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObjectRefByNameFromGroup("alerttype", "trigger", typeID, name, alertTriggerObj)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}

			if err != nil {
				cmd.PrintResult("", err)
			}

			// update the metric object values
			if name != DEFAULT_STRING_FLAG_VALUE {
				alertTriggerObj.Name = name
			}
			if description != DEFAULT_STRING_FLAG_VALUE {
				alertTriggerObj.Description = description
			}
			if autoResolve != -1 {
				alertTriggerObj.AutoResolve = autoResolve
			}
			if metricID != -1 {
				alertTriggerObj.Metric = metricID
			}
			if config != DEFAULT_STRING_FLAG_VALUE {
				alertTriggerObj.Config = config
			}
			if serverGroupID != -1 {
				alertTriggerObj.GroupID = serverGroupID
			}
			if serverID != -1 {
				alertTriggerObj.ServerID = serverID
			}

			onApp = alertTriggerObj.GroupID == 0 && alertTriggerObj.ServerID == 0

			if alertTriggerObj.OnApp != onApp {
				alertTriggerObj.OnApp = onApp
			}

			cmd.PrintResult(cmd.Capi.UpdateTrigger(typeID, alertTriggerObj))
		},
	},
	{
		Name:      "delete",
		UsageLine: `alert trigger delete (--id | --name) (--type | --typeid)`,
		Long: `
Delete a trigger from an alert type group.

The flags for "delete" trigger action are:

Mandatory:
	--id
		Specify the trigger id.
	or
	--name
		Specify the trigger name.
	--typeid
		Specify the alert type id.
	or
	--type
		Specify the alert type name.
`,
		Run: func(cmd *Command, args []string) {
			var id, typeID int64
			var name, Type string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&typeID, "typeid", -1, "Specify the alert type id.")
			cmd.Flag.Int64Var(&id, "id", -1, "Specify the trigger id.")
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Specify the trigger name.")
			cmd.Flag.StringVar(&Type, "type", DEFAULT_STRING_FLAG_VALUE, "Specify the alert type name.")
			cmd.ParseArgs(args)

			// Check the mandatory flags.
			if id == -1 && name == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if typeID == -1 && Type == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}

			var err error
			// get the alert type for the trigger
			var alertTypeObj = &api.AlertType{}
			if typeID == -1 {
				err = cmd.Capi.GetObjectRefByName("alerttype", Type, alertTypeObj)
				if err != nil {
					cmd.PrintResult("", err)
				}
				// if didn't exit due to error...
				typeID = alertTypeObj.ID
			}

			// Get the existing trigger.
			var alertTriggerObj = &api.AlertTrigger{}
			if id == -1 {
				err = cmd.Capi.GetObjectRefByNameFromGroup("alerttype", "trigger", typeID, name, alertTriggerObj)
				if err != nil {
					cmd.PrintResult("", err)
				}
				// if didn't exit due to error...
				id = alertTriggerObj.ID
			}

			cmd.PrintResult(cmd.Capi.DeleteObjectFromGroupByID("alerttype", "trigger", typeID, id))
		},
	},
}
