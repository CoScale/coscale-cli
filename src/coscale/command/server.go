package command

import (
	"coscale/api"
	"os"
)

// ServerObject defines the server command on the CLI.
var ServerObject = NewCommand("server", "server <action> [--<field>='<data>']", ServerActions)

// ServerActions defines the server actions on the CLI.
var ServerActions = []*Command{
	ListCmd("server"),
	GetCmd("server"),
	DeleteCmd(&api.Server{}, "server"),
	{
		Name:      "new",
		UsageLine: "server new (--name) [--description --serverType]",
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
`,
		Run: func(cmd *Command, args []string) {
			var name, description, serverType, source string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the server.")
			cmd.Flag.StringVar(&description, "description", "", "Description for the server.")
			cmd.Flag.StringVar(&serverType, "serverType", "", "Describes the type of server.")
			cmd.Flag.StringVar(&source, "source", "cli", "Deprecated.")
			cmd.ParseArgs(args)

			if name == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			cmd.PrintResult(cmd.Capi.CreateServer(name, description, serverType))
		},
	},
	{
		Name:      "update",
		UsageLine: "server update (--name | --id) [--description --serverType --state]",
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
	--state
	 	"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.
`,
		Run: func(cmd *Command, args []string) {
			var name, description, Type, source, state string
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the server.")
			cmd.Flag.StringVar(&description, "description", DEFAULT_STRING_FLAG_VALUE, "Description for the server.")
			cmd.Flag.StringVar(&Type, "type", DEFAULT_STRING_FLAG_VALUE, "Describes the type of server.")
			cmd.Flag.StringVar(&source, "source", DEFAULT_STRING_FLAG_VALUE, "Deprecated.")
			cmd.Flag.StringVar(&state, "state", DEFAULT_STRING_FLAG_VALUE, `"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.`)
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier.")
			cmd.ParseArgs(args)

			var serverObj = &api.Server{}
			var err error
			if id != -1 {
				err = cmd.Capi.GetObjectRef("server", id, serverObj)
			} else if name != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName("server", name, serverObj)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}

			//update the server object values
			if name != DEFAULT_STRING_FLAG_VALUE {
				serverObj.Name = name
			}
			if description != DEFAULT_STRING_FLAG_VALUE {
				serverObj.Description = description
			}
			if Type != DEFAULT_STRING_FLAG_VALUE {
				serverObj.Type = Type
			}
			if state != DEFAULT_STRING_FLAG_VALUE {
				serverObj.State = state
			}

			cmd.PrintResult(cmd.Capi.UpdateServer(serverObj))
		},
	},
}

// ServerGroupObject defines the server group command on the CLI.
var ServerGroupObject = NewCommand("servergroup", "servergroup <action> [--<field>='<data>']", ServerGroupActions)

// ServerGroupActions defines the server group actions on the CLI.
var ServerGroupActions = []*Command{
	ListCmd("servergroup"),
	GetCmd("servergroup"),
	DeleteCmd(&api.ServerGroup{}, "servergroup"),
	{
		Name:      "new",
		UsageLine: `servergroup new (--name) [--parentId --description --type --state]`,
		Long: `
Create a new CoScale servergroup object.

The flags for new servergroup action are:

Mandatory:
	--name
		Name for the server group.
Optional:
	--parentId
		Optionally set which other ServerGroup will be its parent.
	--description
		Description for the server group.
	--type
		Describes the type of server group.
	--state
		"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.
`,
		Run: func(cmd *Command, args []string) {
			var name, description, Type, state, source string
			var parentID int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the server group.")
			cmd.Flag.Int64Var(&parentID, "parentId", -1, "Optionally set which other ServerGroup will be its parent.")
			cmd.Flag.StringVar(&description, "description", "", "Description for the server group.")
			cmd.Flag.StringVar(&Type, "type", "", "Describes the type of server group.")
			cmd.Flag.StringVar(&state, "state", "", `"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.`)
			cmd.Flag.StringVar(&source, "source", "cli", "Deprecated.")
			cmd.ParseArgs(args)

			if name == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			cmd.PrintResult(cmd.Capi.CreateServerGroup(name, description, Type, state, parentID))
		},
	},
	{
		Name:      "update",
		UsageLine: `servergroup update (--name | --id) [--parentId --description --type --state]`,
		Long: `
Update a CoScale servergroup object.

The flags for update servergroup action are:
The name or id should be specified
	--id
		Unique identifier, if we want to update the name of the servergroup, this become mandatory
	--name
		Name for the server group.
	--parentId
		Optionally set which other ServerGroup will be its parent.
	--description
		Description for the server group.
	--type
		Describes the type of server group.
	--state
		"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.
`,
		Run: func(cmd *Command, args []string) {
			var name, description, Type, source, state string
			var id, parentID int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier.")
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the server group.")
			cmd.Flag.Int64Var(&parentID, "parentId", -1, "Optionally set which other ServerGroup will be its parent.")
			cmd.Flag.StringVar(&description, "description", DEFAULT_STRING_FLAG_VALUE, "Description for the server group.")
			cmd.Flag.StringVar(&Type, "type", DEFAULT_STRING_FLAG_VALUE, "Describes the type of server group.")
			cmd.Flag.StringVar(&state, "state", DEFAULT_STRING_FLAG_VALUE, `"ENABLED": capturing data, "INACTIVE": not capturing data, "DISABLED": not capturing data and not shown on the dashboard.`)
			cmd.Flag.StringVar(&source, "source", DEFAULT_STRING_FLAG_VALUE, "Deprecated.")
			cmd.ParseArgs(args)

			var serverGroupObj = &api.ServerGroup{}
			var err error
			if id != -1 {
				err = cmd.Capi.GetObjectRef("servergroup", id, serverGroupObj)
			} else if name != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName("servergroup", name, serverGroupObj)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}

			//update the server object values
			if name != DEFAULT_STRING_FLAG_VALUE {
				serverGroupObj.Name = name
			}
			if parentID != -1 {
				serverGroupObj.ParentID = parentID
			}
			if description != DEFAULT_STRING_FLAG_VALUE {
				serverGroupObj.Description = description
			}
			if Type != DEFAULT_STRING_FLAG_VALUE {
				serverGroupObj.Type = Type
			}
			if state != DEFAULT_STRING_FLAG_VALUE {
				serverGroupObj.State = state
			}

			cmd.PrintResult(cmd.Capi.UpdateServerGroup(serverGroupObj))
		},
	},
	AddObjToGroupCmd("server", &api.Server{}, &api.ServerGroup{}),
	DeleteObjFromGroupCmd("server", &api.Server{}, &api.ServerGroup{}),
	AddObjToGroupCmd("servergroup", &api.ServerGroup{}, &api.ServerGroup{}),
	DeleteObjFromGroupCmd("servergroup", &api.ServerGroup{}, &api.ServerGroup{}),
}
