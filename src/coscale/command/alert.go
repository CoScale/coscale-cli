package command

import (
	"coscale/api"
	"os"
)

var alertObjectName = "alert"
var AlertObject = NewCommand(alertObjectName, "alert <action> [--<field>='<data>']", AlertActions)

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
			cmd.Flag.StringVar(&filter, "filter", DEFAULT_FLAG_VALUE, "List actions filtered by unresolved or by unacknowledged.")
			cmd.ParseArgs(args)

			switch filter {
			case DEFAULT_FLAG_VALUE:
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
		UsageLine: "alert acknowledge (--id | --name)",
		Long: `
Acknowledge an alert.

The flags for acknowledge alert action are:
Mandatory:
	--id
		The id of the alert.
	or
	--name
		The name of the alert.
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
		UsageLine: "alert resolve (--id | --name)",
		Long: `
Resolve an alert.

The flags for resolve alert action are:
Mandatory:
	--id
		The id of the alert.
	or 
	--name
		The name of the alert.
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
