package command

import (
	"coscale/api"
	"os"
)

var ServerObject = NewCommand("server", "server <action> [--<field>='<data>']", ServerActions)
var ServerActions = []*Command{
	ListCmd("server"),
	GetCmd("server"),
	DeleteCmd("server", &api.Server{}),
	{
		Name:      "new",
		UsageLine: "server new (--name) [--description --serverType --source]",
		Long: `
Create a new CoScale server object.

The flags for new server action are:

Mandatory:
	--name 
		Name for the server.
Optional:
	--description
		Description for the server.
	--serverType
		Describes the type of server.
	--source
		Describes who added the server. Can be chosen by the user. [default: "cli"]
`,
		Run: func(cmd *Command, args []string) {
			var name, description, serverType, source string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_FLAG_VALUE, "Name for the server.")
			cmd.Flag.StringVar(&description, "description", "", "Description for the server.")
			cmd.Flag.StringVar(&serverType, "serverType", "", "Describes the type of server.")
			cmd.Flag.StringVar(&source, "source", "cli", "Describes who added the server.")
			cmd.ParseArgs(args)

			if name == DEFAULT_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			cmd.PrintResult(cmd.Capi.CreateServer(name, description, serverType, source))
		},
	},
	{
		Name:      "update",
		UsageLine: "server update (--name | --id) [--description --serverType --state --source]",
		Long: `
Update a CoScale server object.

The flags for update server action are:
The name or id should be specified
	--id
		Unique identifier, if we want to update the name of the server, this become mandatory
	--name 
		specify the name of the server.
	--description
		Description for the server.
	--serverType
		Describes the type of server.
	--source
		Describes who added the server. Can be chosen by the user. [default: "cli"]
	--state
	 	"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.
`,
		Run: func(cmd *Command, args []string) {
			var name, description, Type, source, state string
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_FLAG_VALUE, "Name for the server.")
			cmd.Flag.StringVar(&description, "description", DEFAULT_FLAG_VALUE, "Description for the server.")
			cmd.Flag.StringVar(&Type, "type", DEFAULT_FLAG_VALUE, "Describes the type of server.")
			cmd.Flag.StringVar(&source, "source", DEFAULT_FLAG_VALUE, "Describes who added the server.")
			cmd.Flag.StringVar(&state, "state", DEFAULT_FLAG_VALUE, `"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.`)
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier")
			cmd.ParseArgs(args)

			var serverObj = &api.Server{}
			var err error
			if id != -1 {
				err = cmd.Capi.GetObjectRef("server", id, serverObj)
			} else if name != DEFAULT_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName("server", name, serverObj)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}

			//update the server object values
			if name != DEFAULT_FLAG_VALUE {
				serverObj.Name = name
			}
			if description != DEFAULT_FLAG_VALUE {
				serverObj.Description = description
			}
			if Type != DEFAULT_FLAG_VALUE {
				serverObj.Type = Type
			}
			if source != DEFAULT_FLAG_VALUE {
				serverObj.Source = source
			}
			if state != DEFAULT_FLAG_VALUE {
				serverObj.State = state
			}

			cmd.PrintResult(cmd.Capi.UpdateServer(serverObj))
		},
	},
}

var ServerGroupObject = NewCommand("servergroup", "servergroup <action> [--<field>='<data>']", ServerGroupActions)
var ServerGroupActions = []*Command{
	ListCmd("servergroup"),
	GetCmd("servergroup"),
	DeleteCmd("servergroup", &api.ServerGroup{}),
	{
		Name:      "new",
		UsageLine: `servergroup new (--name) [--description --type --state --source]`,
		Long: `
Create a new CoScale servergroup object.

The flags for new servergroup action are:

Mandatory:
	--name 
		Name for the server group.
Optional:
	--description
		Description for the server group.
	--type
		Describes the type of server group.
	--state
		"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.	
	--source
		Describes who added the server group. Can be chosen by the user. [default: "cli"]
`,
		Run: func(cmd *Command, args []string) {
			var name, description, Type, state, source string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_FLAG_VALUE, "Name for the server group.")
			cmd.Flag.StringVar(&description, "description", "", "Description for the server group.")
			cmd.Flag.StringVar(&Type, "type", "", "Describes the type of server group.")
			cmd.Flag.StringVar(&state, "state", "", `"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.`)
			cmd.Flag.StringVar(&source, "source", "cli", "Describes who added the server group.")
			cmd.ParseArgs(args)

			if name == DEFAULT_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			cmd.PrintResult(cmd.Capi.CreateServerGroup(name, description, Type, state, source))
		},
	},
	{
		Name:      "update",
		UsageLine: `servergroup update (--name | --id) [--description --type --state --source]`,
		Long: `
Update a CoScale servergroup object.

The flags for update servergroup action are:
The name or id should be specified
	--id
		Unique identifier, if we want to update the name of the servergroup, this become mandatory
	--name 
		Name for the server group.
	--description
		Description for the server group.
	--type
		Describes the type of server group.
	--state
		"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.	
	--source
		Describes who added the server group. Can be chosen by the user. [default: "cli"]
`,
		Run: func(cmd *Command, args []string) {
			var name, description, Type, source, state string
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier.")
			cmd.Flag.StringVar(&name, "name", DEFAULT_FLAG_VALUE, "Name for the server group.")
			cmd.Flag.StringVar(&description, "description", DEFAULT_FLAG_VALUE, "Description for the server group.")
			cmd.Flag.StringVar(&Type, "type", DEFAULT_FLAG_VALUE, "Describes the type of server group.")
			cmd.Flag.StringVar(&state, "state", DEFAULT_FLAG_VALUE, `"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.`)
			cmd.Flag.StringVar(&source, "source", DEFAULT_FLAG_VALUE, "Describes who added the server. Can be chosen by the user.")
			cmd.ParseArgs(args)

			var serverGroupObj = &api.ServerGroup{}
			var err error
			if id != -1 {
				err = cmd.Capi.GetObjectRef("servergroup", id, serverGroupObj)
			} else if name != DEFAULT_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName("servergroup", name, serverGroupObj)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}

			//update the server object values
			if name != DEFAULT_FLAG_VALUE {
				serverGroupObj.Name = name
			}
			if description != DEFAULT_FLAG_VALUE {
				serverGroupObj.Description = description
			}
			if Type != DEFAULT_FLAG_VALUE {
				serverGroupObj.Type = Type
			}
			if source != DEFAULT_FLAG_VALUE {
				serverGroupObj.Source = source
			}
			if state != DEFAULT_FLAG_VALUE {
				serverGroupObj.State = state
			}

			cmd.PrintResult(cmd.Capi.UpdateServerGroup(serverGroupObj))
		},
	},
	AddObjToGroupCmd("server", &api.Server{}, &api.ServerGroup{}),
	DeleteObjFromGroupCmd("server", &api.Server{}, &api.ServerGroup{}),
}
