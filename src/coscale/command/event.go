package command

import (
	"coscale/api"
	"os"
)

var eventObjectName = "event"
var EventObject = NewCommand(eventObjectName, "event <action> [--<field>='<data>']", EventActions)

var EventActions = []*Command{
	ListCmd(eventObjectName),
	GetCmd(eventObjectName),
	DeleteCmd(eventObjectName, &api.Event{}),
	{
		Name:      "new",
		UsageLine: "event new (--name) [--description --attributeDescriptions --source]",
		Long: `
Create new event category.

The flags for new event action are:

Mandatory:
	--name 
		specify name of the event.
Optional:
	--description
		specify the description of the event.
	--attributeDescriptions
		JSON string describing what items the "attribute" of an EventData instance belonging to this Event must have.  [default: "[]"]
	--source
		Describes who added the event. Can be chosen by the user. [default: "cli"]
`,
		Run: func(cmd *Command, args []string) {
			var name, description, attributeDescriptions, source string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_FLAG_VALUE, "specify the event name of the event")
			cmd.Flag.StringVar(&description, "description", "", "specify the description of the event")
			cmd.Flag.StringVar(&attributeDescriptions, "attributeDescriptions", "[]", "")
			cmd.Flag.StringVar(&source, "source", "cli", "Describes who added the event")
			cmd.ParseArgs(args)

			if name == DEFAULT_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			cmd.PrintResult(cmd.Capi.CreateEvent(name, description, attributeDescriptions, source, ""))
		},
	},
	{
		Name:      "update",
		UsageLine: "event update (--name | --id) [--description --attributeDescriptions --source]",
		Long: `
Update a CoScale event object.

The flags for update event action are:
The name or id should be specified
	--id
		Unique identifier, if we want to update the name of the event, this become mandatory
	--name 
		specify the event name of the event.
	--description
		specify the description of the event.
	--attributeDescriptions
		JSON string describing what items the "attribute" of an EventData instance belonging to this Event must have.
	--source
		Describes who added the event. Can be chosen by the user. [default: "cli"]
`,
		Run: func(cmd *Command, args []string) {
			var name, description, attributeDescriptions, source string
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_FLAG_VALUE, "specify the event name of the event")
			cmd.Flag.StringVar(&description, "description", DEFAULT_FLAG_VALUE, "specify the description of the event")
			cmd.Flag.StringVar(&attributeDescriptions, "attributeDescriptions", DEFAULT_FLAG_VALUE, "")
			cmd.Flag.StringVar(&source, "source", DEFAULT_FLAG_VALUE, "Describes who added the event")
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier")
			cmd.ParseArgs(args)

			var eventObj = &api.Event{}
			var err error
			if id != -1 {
				err = cmd.Capi.GetObjectRef("event", id, eventObj)
			} else if name != DEFAULT_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName("event", name, eventObj)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}
			//update the event object values
			if name != DEFAULT_FLAG_VALUE {
				eventObj.Name = name
			}
			if description != DEFAULT_FLAG_VALUE {
				eventObj.Description = description
			}
			if attributeDescriptions != DEFAULT_FLAG_VALUE {
				eventObj.AttributeDescriptions = attributeDescriptions
			}
			if source != DEFAULT_FLAG_VALUE {
				eventObj.Source = source
			}

			cmd.PrintResult(cmd.Capi.UpdateEvent(eventObj))
		},
	},
	{
		Name:      "data",
		UsageLine: "event data (--name --id --message --subject) [--attribute --timestamp --stopTime]",
		Long: `
Insert event data.

The flags for data event action are:
Mandatory:
	--name 
		specify the event name.
	--id
		specify the event id.
	Only one from id/name is neccessary.
		
	--message
		The message for the event data.	
	--subject
		The subject for the event data. The subject is structured as follows:
		s<serverId> for a server, g<servergroupId> for a server group, a for the application.
Optional:	
	--attribute
		JSON String detailing the progress of the event.
	--timestamp
		Timestamp in seconds ago. [default: 0]
	--stopTime
		The time at which the EventData stopped in seconds ago. [default: 0]
		`,
		Run: func(cmd *Command, args []string) {
			var id, timestamp, stopTime int64
			var name, message, subject, attribute string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_FLAG_VALUE, "event name")
			cmd.Flag.StringVar(&message, "message", DEFAULT_FLAG_VALUE, "message for the event data")
			cmd.Flag.StringVar(&subject, "subject", DEFAULT_FLAG_VALUE, "subject for the event data")
			cmd.Flag.StringVar(&attribute, "attribute", "{}", "JSON String detailing the progress of the event")
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier")
			cmd.Flag.Int64Var(&timestamp, "timestamp", 0, "Timestamp in seconds ago")
			cmd.Flag.Int64Var(&stopTime, "stopTime", 0, "The time at which the EventData stopped in seconds ago")
			cmd.ParseArgs(args)

			var eventObj = &api.Event{}
			var err error
			//check for mandatory args
			flagErr := message == DEFAULT_FLAG_VALUE || subject == DEFAULT_FLAG_VALUE
			if !flagErr {
				if id != -1 {
					err = cmd.Capi.GetObjectRef("event", id, eventObj)
				} else if name != DEFAULT_FLAG_VALUE {
					err = cmd.Capi.GetObejctRefByName("event", name, eventObj)
				} else {
					flagErr = true
				}
			}
			if flagErr {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}
			// be sure that the time values are negative
			if timestamp > 0 {
				timestamp = -timestamp
			}
			if stopTime > 0 {
				stopTime = -stopTime
			}
			cmd.PrintResult(cmd.Capi.InsertEventData(eventObj.ID, message, subject, attribute, timestamp, stopTime))
		},
	},
}
