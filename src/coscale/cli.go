package main

import (
	"coscale/command"
	"flag"
	"os"
)

const FIRST_RUN_SUCCESS int = 67

//main command of Coscale coscale-cli
func main() {

	var subCommands = []*command.Command{
		command.EventObject,
		command.ServerObject,
		command.ServerGroupObject,
		command.MetricObject,
		command.MetricGroupObject,
		command.DataObject,
		command.AlertObject,
		command.CheckObject,
	}
	var usage = os.Args[0] + ` <object> <action> [--<field>='<data>']`
	var app = command.NewCommand(os.Args[0], usage, subCommands)
	flag.Usage = func() { app.PrintUsage() }
	var firstRun bool
	flag.BoolVar(&firstRun, "first-run", false, "The first run of an updated agent: checks if everything is working properly.")
	flag.Parse()
	if firstRun {
		os.Exit(FIRST_RUN_SUCCESS)
	}
	app.Run(app, flag.Args())
}
