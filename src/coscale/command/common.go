package command

import (
	"coscale/api"
	"fmt"
	"os"
)

// parseParams will be used in cases that we want
// to have a different command name than the name of the resource.
func parseParams(params []string) (string, string) {
	switch len(params) {
	case 0:
		return "", ""
	case 1:
		return params[0], params[0]
	default:
		return params[0], params[1]
	}
}

// ListCmd is used to create a command that will list all the objects with a given type.
func ListCmd(params ...string) *Command {
	objectName, cmdName := parseParams(params)
	return &Command{
		Name:      "list",
		UsageLine: fmt.Sprintf("%s list", cmdName),
		Long: fmt.Sprintf(`
Get all %[1]ss from CoScale Api.
`, objectName),
		Run: func(cmd *Command, args []string) {
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.ParseArgs(args)
			cmd.PrintResult(cmd.Capi.GetObjects(objectName))
		},
	}
}

// GetCmd is used to create a command that will get an object with a given type.
func GetCmd(params ...string) *Command {
	objectName, cmdName := parseParams(params)
	return &Command{
		Name:      "get",
		UsageLine: fmt.Sprintf("%s get (--id | --name)", cmdName),
		Long: fmt.Sprintf(`
Get a CoScale %[1]s object by id or by name.

The flags for %[2]s get action are:
Only one of them is necessary to be specified
	--name 
		specify the %[1]s name.
	--id
		specify the %[1]s id.
`, objectName, cmdName),
		Run: func(cmd *Command, args []string) {
			var name string
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the object.")
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier.")
			cmd.ParseArgs(args)

			if id != -1 {
				cmd.PrintResult(cmd.Capi.GetObject(objectName, id))
			} else if name != DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintResult(cmd.Capi.GetObjectByName(objectName, name))
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
		},
	}
}

// DeleteCmd is used to create a command that will delete an object with a given type.
func DeleteCmd(object api.Object, params ...string) *Command {
	objectName, cmdName := parseParams(params)
	return &Command{
		Name:      "delete",
		UsageLine: fmt.Sprintf("%s delete (--name | --id)", cmdName),
		Long: fmt.Sprintf(`
Delete a %[1]s by the name or id.

The flags for %[2]s delete action are:
Only one of them is necessary to be specified
	--name 
		specify the %[1]s name.
	--id
		specify the %[1]s id.
`, objectName, cmdName),
		Run: func(cmd *Command, args []string) {
			var name string
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Name for the object.")
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier.")
			cmd.ParseArgs(args)

			var err error
			if id != -1 {
				err = cmd.Capi.GetObjectRef(objectName, id, object)
			} else if name != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName(objectName, name, object)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}
			cmd.PrintResult("", cmd.Capi.DeleteObject(objectName, &object))
		},
	}
}

// AddObjToGroupCmd is used create a command that will add an object to a group e.g. metric to metricgroup.
func AddObjToGroupCmd(objectName string, object api.Object, group api.Object) *Command {
	return &Command{
		Name:      fmt.Sprintf("add%s", capitalize(objectName)),
		UsageLine: fmt.Sprintf(`%[1]sgroup add%[2]s (--id%[2]s | --name%[2]s) (--idGroup | --nameGroup)`, objectName, capitalize(objectName)),
		Long: fmt.Sprintf(`
Add a existing %[1]s to a %[1]s group.

The flags for "add%[2]s" %[1]sgroup action are:

Mandatory:
	--id%[2]s 
		specify the %[1]s id.
	or	
	--name%[2]s
		specify the %[1]s name.
	--idGroup 
		specify the group id.
	or	
	--nameGroup
		specify the group name.	
`, objectName, capitalize(objectName)),
		Run: func(cmd *Command, args []string) {
			var idGroup, idObject int64
			var nameGroup, nameObject string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&idGroup, "idGroup", -1, "")
			cmd.Flag.Int64Var(&idObject, fmt.Sprintf("id%s", capitalize(objectName)), -1, "")
			cmd.Flag.StringVar(&nameGroup, "nameGroup", DEFAULT_STRING_FLAG_VALUE, "")
			cmd.Flag.StringVar(&nameObject, fmt.Sprintf("name%s", capitalize(objectName)), DEFAULT_STRING_FLAG_VALUE, "")
			cmd.ParseArgs(args)

			var err error
			if idObject != -1 {
				err = cmd.Capi.GetObjectRef(objectName, idObject, object)
			} else if nameObject != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName(objectName, nameObject, object)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}

			if idGroup != -1 {
				err = cmd.Capi.GetObjectRef(fmt.Sprintf("%sgroup", objectName), idGroup, group)
			} else if nameObject != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName(fmt.Sprintf("%sgroup", objectName), nameGroup, group)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}
			cmd.PrintResult("", cmd.Capi.AddObjectToGroup(objectName, object, group))
		},
	}
}

// DeleteObjFromGroupCmd is used create a command that will delete an object from a group e.g. metric to metricgroup.
func DeleteObjFromGroupCmd(objectName string, object api.Object, group api.Object) *Command {
	return &Command{
		Name:      fmt.Sprintf("delete%s", capitalize(objectName)),
		UsageLine: fmt.Sprintf(`%[1]sgroup delete%[2]s (--id%[2]s | --name%[2]s) (--idGroup | --nameGroup)`, objectName, capitalize(objectName)),
		Long: fmt.Sprintf(`
Delete a %[1]s from a %[1]s group.

The flags for "delete%[2]s" %[1]sgroup action are:

Mandatory:
	--id%[2]s 
		specify the %[1]s id.
	or	
	--name%[2]s
		specify the %[1]s name.
	--idGroup 
		specify the group id.
	or	
	--nameGroup
		specify the group name.	
`, objectName, capitalize(objectName)),
		Run: func(cmd *Command, args []string) {
			var idGroup, idObject int64
			var nameGroup, nameObject string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&idGroup, "idGroup", -1, "")
			cmd.Flag.Int64Var(&idObject, fmt.Sprintf("id%s", capitalize(objectName)), -1, "")
			cmd.Flag.StringVar(&nameGroup, "nameGroup", DEFAULT_STRING_FLAG_VALUE, "")
			cmd.Flag.StringVar(&nameObject, fmt.Sprintf("name%s", capitalize(objectName)), DEFAULT_STRING_FLAG_VALUE, "")
			cmd.ParseArgs(args)

			var err error
			if idObject != -1 {
				err = cmd.Capi.GetObjectRef(objectName, idObject, object)
			} else if nameObject != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName(objectName, nameObject, object)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}
			if idGroup != -1 {
				err = cmd.Capi.GetObjectRef(fmt.Sprintf("%sgroup", objectName), idGroup, group)
			} else if nameObject != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName(fmt.Sprintf("%sgroup", objectName), nameGroup, group)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}
			cmd.PrintResult("", cmd.Capi.DeleteObjectFromGroup(objectName, object, group))
		},
	}
}
