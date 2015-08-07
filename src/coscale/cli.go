package main

import (
	"coscale/command"
	"flag"
	"os"
)

//main command of Coscale cli-cmd
func main() {

	var subCommands = []*command.Command{
		command.EventObject,
		command.ServerObject,
		command.ServerGroupObject,
		command.MetricObject,
		command.MetricGroupObject,
		command.DataObject,
		command.AlertObject,
	}
	var usage = os.Args[0] + ` <object> <action> [--<field>='<data>']`
	var app = command.NewCommand(os.Args[0], usage, subCommands)
	flag.Usage = func() { app.PrintUsage() }
	flag.Parse()
	app.Run(app, flag.Args())
}
