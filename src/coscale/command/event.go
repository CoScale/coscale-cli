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
	DeleteCmd(&api.Event{}, eventObjectName),
	{
		Name:      "new",
		UsageLine: "event new (--name) [--description --attributeDescriptions]",
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
`,
		Run: func(cmd *Command, args []string) {
			var name, eventType, description, attributeDescriptions, source string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Specify the name of the event.")
			cmd.Flag.StringVar(&eventType, "type", "", "Specify the type of the event.")
			cmd.Flag.StringVar(&description, "description", "", "Specify the description of the event.")
			cmd.Flag.StringVar(&attributeDescriptions, "attributeDescriptions", "[]", "JSON string describing what items the attribute.")
			cmd.Flag.StringVar(&source, "source", "cli", "Deprecated.")
			cmd.ParseArgs(args)

			if name == DEFAULT_STRING_FLAG_VALUE {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			cmd.PrintResult(cmd.Capi.CreateEvent(name, description, attributeDescriptions, eventType))
		},
	},
	{
		Name:      "update",
		UsageLine: "event update (--name | --id) [--description --attributeDescriptions]",
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
`,
		Run: func(cmd *Command, args []string) {
			var name, eventType, description, attributeDescriptions, source string
			var id int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Specify the name of the event.")
			cmd.Flag.StringVar(&eventType, "type", DEFAULT_STRING_FLAG_VALUE, "Specify the type of the event.")
			cmd.Flag.StringVar(&description, "description", DEFAULT_STRING_FLAG_VALUE, "Specify the description of the event.")
			cmd.Flag.StringVar(&attributeDescriptions, "attributeDescriptions", DEFAULT_STRING_FLAG_VALUE, "JSON string describing what items the attribute.")
			cmd.Flag.StringVar(&source, "source", DEFAULT_STRING_FLAG_VALUE, "Deprecated.")
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier")
			cmd.ParseArgs(args)

			var eventObj = &api.Event{}
			var err error
			if id != -1 {
				err = cmd.Capi.GetObjectRef("event", id, eventObj)
			} else if name != DEFAULT_STRING_FLAG_VALUE {
				err = cmd.Capi.GetObejctRefByName("event", name, eventObj)
			} else {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			if err != nil {
				cmd.PrintResult("", err)
			}
			//update the event object values
			if name != DEFAULT_STRING_FLAG_VALUE {
				eventObj.Name = name
			}
			if description != DEFAULT_STRING_FLAG_VALUE {
				eventObj.Description = description
			}
			if attributeDescriptions != DEFAULT_STRING_FLAG_VALUE {
				eventObj.AttributeDescriptions = attributeDescriptions
			}

			if eventType != DEFAULT_STRING_FLAG_VALUE {
				eventObj.Type = eventType
			}

			cmd.PrintResult(cmd.Capi.UpdateEvent(eventObj))
		},
	},
	{
		Name:      "data",
		UsageLine: "event data (--name --id --message --subject) [--attribute --timestamp --stopTime]",
		Long: `
		Warning! 'event data' is deprecated and will be removed in the future.
Please use 'event newdata' instead.
		`,
		Deprecated: true,
		Run: func(cmd *Command, args []string) {
			var id, timestamp, stopTime int64
			var name, message, subject, attribute string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Event name.")
			cmd.Flag.StringVar(&message, "message", DEFAULT_STRING_FLAG_VALUE, "Message for the event data.")
			cmd.Flag.StringVar(&subject, "subject", DEFAULT_STRING_FLAG_VALUE, "Subject for the event data.")
			cmd.Flag.StringVar(&attribute, "attribute", "{}", "JSON String detailing the progress of the event.")
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier.")
			cmd.Flag.Int64Var(&timestamp, "timestamp", 0, "Timestamp in seconds ago.")
			cmd.Flag.Int64Var(&stopTime, "stopTime", DEFAULT_INT64_FLAG_VALUE, "The time at which the EventData stopped in seconds ago.")
			cmd.ParseArgs(args)

			var eventObj = &api.Event{}
			var err error
			//check for mandatory args
			flagErr := message == DEFAULT_STRING_FLAG_VALUE || subject == DEFAULT_STRING_FLAG_VALUE
			if !flagErr {
				if id != -1 {
					err = cmd.Capi.GetObjectRef("event", id, eventObj)
				} else if name != DEFAULT_STRING_FLAG_VALUE {
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
			cmd.PrintResult(cmd.Capi.InsertEventData(eventObj.ID, message, subject, attribute, timestamp, stopTime))
		},
	},
	{
		Name:      "newdata",
		UsageLine: "event newdata (--name --id --message --subject) [--attribute --timestamp --stopTime]",
		Long: `
Insert event data.

The flags for newdata event action are:
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
		Timestamp in seconds ago(negative values) or unix timestamp(positive values). [default: 0]
	--stopTime
		The time at which the EventData stopped in seconds ago(negative values) or unix timestamp(positive values).
		`,
		Run: func(cmd *Command, args []string) {
			var id, timestamp, stopTime int64
			var name, message, subject, attribute string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Event name.")
			cmd.Flag.StringVar(&message, "message", DEFAULT_STRING_FLAG_VALUE, "Message for the event data.")
			cmd.Flag.StringVar(&subject, "subject", DEFAULT_STRING_FLAG_VALUE, "Subject for the event data.")
			cmd.Flag.StringVar(&attribute, "attribute", "{}", "JSON String detailing the progress of the event.")
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier.")
			cmd.Flag.Int64Var(&timestamp, "timestamp", 0, "Timestamp in seconds ago.")
			cmd.Flag.Int64Var(&stopTime, "stopTime", DEFAULT_INT64_FLAG_VALUE, "The time at which the EventData stopped in seconds ago.")
			cmd.ParseArgs(args)

			var eventObj = &api.Event{}
			var err error
			//check for mandatory args
			flagErr := message == DEFAULT_STRING_FLAG_VALUE || subject == DEFAULT_STRING_FLAG_VALUE
			if !flagErr {
				if id != -1 {
					err = cmd.Capi.GetObjectRef("event", id, eventObj)
				} else if name != DEFAULT_STRING_FLAG_VALUE {
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
			cmd.PrintResult(cmd.Capi.InsertEventData(eventObj.ID, message, subject, attribute, timestamp, stopTime))
		},
	},
	{
		Name:      "updatedata",
		UsageLine: "event updatedata (--name --id --dataid) [ --message --subject --attribute --timestamp --stopTime]",
		Long: `
Update event data.

The flags for updatedata event action are:
Mandatory:
	--name 
		specify the event name.
	--id
		specify the event id.
	Only one from id/name is neccessary.
	
	--dataid
		specify the unique id of the event data.
Optional:
	--message
		The message for the event data.	
	--subject
		The subject for the event data. The subject is structured as follows:
		s<serverId> for a server, g<servergroupId> for a server group, a for the application.
	--attribute
		JSON String detailing the progress of the event.
	--timestamp
		Timestamp in seconds ago(negative values) or unix timestamp(positive values).
	--stopTime
		The time at which the EventData stopped in seconds ago(negative values) or unix timestamp(positive values).
		`,
		Run: func(cmd *Command, args []string) {
			var id, dataid, timestamp, stopTime int64
			var name, message, subject, attribute string
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.StringVar(&name, "name", DEFAULT_STRING_FLAG_VALUE, "Event name.")
			cmd.Flag.StringVar(&message, "message", DEFAULT_STRING_FLAG_VALUE, "Message for the event data.")
			cmd.Flag.StringVar(&subject, "subject", DEFAULT_STRING_FLAG_VALUE, "Subject for the event data.")
			cmd.Flag.StringVar(&attribute, "attribute", "{}", "JSON String detailing the progress of the event.")
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier.")
			cmd.Flag.Int64Var(&dataid, "dataid", -1, "Unique identifier of the event data.")
			cmd.Flag.Int64Var(&timestamp, "timestamp", DEFAULT_INT64_FLAG_VALUE, "Timestamp in seconds ago.")
			cmd.Flag.Int64Var(&stopTime, "stopTime", DEFAULT_INT64_FLAG_VALUE, "The time at which the EventData stopped in seconds ago.")
			cmd.ParseArgs(args)

			var eventObj = &api.Event{}
			var err error
			//check for mandatory args
			flagErr := dataid == -1
			if !flagErr {
				if id != -1 {
					err = cmd.Capi.GetObjectRef("event", id, eventObj)
				} else if name != DEFAULT_STRING_FLAG_VALUE {
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

			// get the existing eventdata for the eventid
			var eventDataObj = &api.EventData{}
			if err = cmd.Capi.GetEventData(eventObj.ID, dataid, eventDataObj); err != nil {
				cmd.PrintResult("", err)
			}
			//update the eventdata object values if are not the default values
			if message != DEFAULT_STRING_FLAG_VALUE {
				eventDataObj.Message = message
			}
			if subject != DEFAULT_STRING_FLAG_VALUE {
				eventDataObj.Subject = subject
			}
			if attribute != `{}` {
				eventDataObj.Attribute = attribute
			}
			if timestamp != DEFAULT_INT64_FLAG_VALUE {
				eventDataObj.Timestamp = timestamp
			}
			eventDataObj.Stoptime = stopTime
			cmd.PrintResult(cmd.Capi.UpdateEventData(eventObj.ID, dataid, eventDataObj))
		},
	},
	{
		Name:      "deletedata",
		UsageLine: "event deletedata (--id --dataid)",
		Long: `
Delete a eventdata entry.

The flags for event deletedata action are:
Mandatory:
	--id 
		specify the event id.
	--dataid
		specify the unique id of the event data.
`,
		Run: func(cmd *Command, args []string) {
			var id, dataid int64
			cmd.Flag.Usage = func() { cmd.PrintUsage() }
			cmd.Flag.Int64Var(&id, "id", -1, "Unique identifier.")
			cmd.Flag.Int64Var(&dataid, "dataid", -1, "Specify the unique id of the event data.")
			cmd.ParseArgs(args)

			// check the args
			if id == -1 || dataid == -1 {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}
			cmd.PrintResult("", cmd.Capi.DeleteEventData(id, dataid))
		},
	},
}
