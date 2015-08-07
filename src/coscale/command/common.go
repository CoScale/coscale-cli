package command

import (
	"coscale/api"
	"fmt"
	"os"
)

func ListCmd(objectName string) *Command {
	return &Command{
		Name:      "list",
		UsageLine: fmt.Sprintf("%s list", objectName),
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

func GetCmd(objectName string) *Command {
	return &Command{
		Name:      "get",
		UsageLine: fmt.Sprintf("%s get (--id | --name)", objectName),
		Long: fmt.Sprintf(`
Get a CoScale %[1]s object by id or by name.

The flags for %[1]s get action are:
Only one of them is necessary to be specified
	--name 
		specify the %[1]s name.
	--id
		specify the %[1]s id.
`, objectName),
		Run: func(cmd *Command, args []string) {
			var name string
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_FLAG_VALUE, "Name for the object")
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier")
			cmd.ParseArgs(args)

			if id != -1 {
				cmd.PrintResult(cmd.Capi.GetObject(objectName, id))
			} else if name != DEFAULT_FLAG_VALUE {
				cmd.PrintResult(cmd.Capi.GetObjectByName(objectName, name))
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
		},
	}
}

func DeleteCmd(objectName string, object api.Object) *Command {
	return &Command{
		Name:      "delete",
		UsageLine: fmt.Sprintf("%s delete (--name | --id)", objectName),
		Long: fmt.Sprintf(`
Delete a %[1]s by the name or id.

The flags for %[1]s delete action are:
Only one of them is necessary to be specified
	--name 
		specify the %[1]s name.
	--id
		specify the %[1]s id.
`, objectName),
		Run: func(cmd *Command, args []string) {
			var name string
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_FLAG_VALUE, "Name for the object")
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier")
			cmd.ParseArgs(args)

			var err error
			if id != -1 {
				err = cmd.Capi.GetObjectRef(objectName, id, object)
			} else if name != DEFAULT_FLAG_VALUE {
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
			cmd.Flag.StringVar(&nameGroup, "nameGroup", DEFAULT_FLAG_VALUE, "")
			cmd.Flag.StringVar(&nameObject, fmt.Sprintf("name%s", capitalize(objectName)), DEFAULT_FLAG_VALUE, "")
			cmd.ParseArgs(args)

			var err error
			if idObject != -1 {
				err = cmd.Capi.GetObjectRef(objectName, idObject, object)
			} else if nameObject != DEFAULT_FLAG_VALUE {
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
			} else if nameObject != DEFAULT_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName(fmt.Sprintf("%sgroup", objectName), nameObject, group)
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
			cmd.Flag.StringVar(&nameGroup, "nameGroup", DEFAULT_FLAG_VALUE, "")
			cmd.Flag.StringVar(&nameObject, fmt.Sprintf("name%s", capitalize(objectName)), DEFAULT_FLAG_VALUE, "")
			cmd.ParseArgs(args)

			var err error
			if idObject != -1 {
				err = cmd.Capi.GetObjectRef(objectName, idObject, object)
			} else if nameObject != DEFAULT_FLAG_VALUE {
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
			} else if nameObject != DEFAULT_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName(fmt.Sprintf("%sgroup", objectName), nameObject, group)
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
