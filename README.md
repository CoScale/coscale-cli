[![Build Status](https://travis-ci.org/CoScale/coscale-cli.svg?branch=master)](https://travis-ci.org/CoScale/coscale-cli) [![Go Report Card](https://goreportcard.com/badge/github.com/coscale/coscale-cli)](https://goreportcard.com/report/github.com/coscale/coscale-cli)

# CoScale command line interface

Communicate with the CoScale API using this convenient Command Line Interface.

The CLI allows you to automate the following actions:

* Manage **Servers**, **Metrics**, **Events** and **Alerts**
* Push and retrieve data for **Metrics** and **Events**

## Installation

The easiest way to use the CLI is by using the Docker image that is available on the Docker hub.

```
docker run coscale/cli [your command] --app-id=[application_id] --access-token=[accesstoken]
```

In bash you can create a shortcut by putting the following function in ~/.bashrc:

```bash
function coscale-cli {
    export ARGS=''; for ARG in "$@"; do if [[ "$ARG" != "${ARG%[[:space:]]*}" ]]; then ARG=\'$ARG\'; fi; ARGS="$ARGS $ARG"; done
    bash -c "docker run coscale/cli ${ARGS} --app-id=[application_id] --access-token=[access_token]"
}
```

Don't forget to fill in your *application id* and *access token* as provided on the Access Token page in the CoScale UI. Restart you bash terminal, you can now use **coscale-cli** without having to provide the applicaiton id and access token every time.

## Usage

```
Usage:
coscale-cli <object> <action> [--<field>='<data>']

Actions for command "coscale-cli":

event
        event <action> [--<field>='<data>']
server
        server <action> [--<field>='<data>']
servergroup
        servergroup <action> [--<field>='<data>']
metric
        metric <action> [--<field>='<data>']
metricgroup
        metricgroup <action> [--<field>='<data>']
data
        data <action> [--<field>='<data>']
alert
        alert <action> [--<field>='<data>']
```

## Examples

### Event Examples

#### Create a new event category

Create an event category called "Releases"

```
coscale-cli event new --name "Releases"
```

This returns an event category object

```
{
 "id": 53,
 "version": 1,
 "state": "ENABLED",
 "name": "Releases",
 "description": "",
 "attributeDescriptions": "[]",
 "type": "",
 "source": "cli",
 "icon": null
}
```

#### Add a new event

Add a CoScale event to the "Releases" event category

```
coscale-cli event newdata --name "Releases" --message "V2.3.4" --subject "a"
```

This returns the event object

```
{
 "gid": "AVwbfl-1EjRUeON0eWk5",
 "id": 16215,
 "applicationId": 21,
 "timestamp": 1495109885,
 "stopTime": null,
 "updateTime": 1495109885,
 "message": "V2.3.4",
 "subject": "a",
 "attribute": "{}",
 "version": 1,
 "eventId": 53
}
```


### Server Examples

#### Add a new server

Add a server to CoScale that that will not run a CoScale agent, but can be used to push custom events and metrics to.

```
coscale-cli server new --name "cron server" --description "Server that that runs cron jobs"
```

This returns the server object

```
{
 "id": 341,
 "version": 1,
 "state": "ENABLED",
 "name": "cron server",
 "description": "Server that that runs cron jobs",
 "type": "",
 "source": "cli",
 "startTime": null,
 "stopTime": null
}
```

#### Delete a server

Delete the cron server.

```
coscale-cli server delete --name "cron server"
```

This does not return any output when successful, the exit code should be 0.


### Alert Examples

#### View all unresolved alerts.

```
coscale-cli alert list --filter unresolved
```

This returns a list of unresolved alerts.

```
[
 {
  "id": 24,
  "version": 1,
  "created": 1495015602,
  "lastOccurence": 1495110060,
  "occurrences": 195,
  "sent": 1495015637,
  "backup": null,
  "escalation": null,
  "acknowledged": null,
  "acknowledgedBy": null,
  "resolved": null,
  "resolvedBy": null,
  "onApp": false,
  "config": "agentTimeout() > 300",
  "acknowledgeMailSent": null,
  "autoresolveMailSent": null,
  "anomalyMessage": null,
  "dimensionSpec": "[]",
  "thirdparty": null,
  "metricId": null,
  "serverId": 271,
  "groupId": null,
  "aggregatedIntoId": null
 }
]
```

#### Resolve an alert from the list.

Use the alert id to resolve the unresolved alert. In this example we will use the alert id from the previous *alert list* output.

```
coscale-cli alert resolve --id 24
```


### Metric Examples

#### Push custom data for an application metric.

Create an application metric

```
coscale-cli metric new --name "Transaction value" --dataType "DOUBLE" --subject "APPLICATION" --unit "$"
```

This returns the metric object (we will need the id for inserting data)

```
{
 "id": 675,
 "version": 1,
 "state": "ENABLED",
 "name": "Transaction value",
 "description": "",
 "unit": "$",
 "period": 60,
 "source": "cli",
 "dataType": "DOUBLE",
 "subject": "APPLICATION"
}
```

Push data for the application metric. In this example we will push data for the 'Transaction value' metric (id 675) at timestamp 1495108650 with value 1.23:

```
coscale-cli data insert --data="M675:A:1495108650:1.23"
```

#### Push custom data for a server metric.

Get the id of the server

```
coscale-cli server get --name 'myserver'
```

This returns the matching server objects

```
[
 {
  "id": 34,
  "version": 3,
  "state": "ENABLED",
  "name": "myserver",
  "description": "Created by agent for 'myserver'",
  "type": "",
  "source": "CoScale Agent",
  "startTime": 1485958383,
  "stopTime": null
 }
]
```

Create a server metric

```
coscale-cli metric new --name "Core temperature" --dataType "DOUBLE" --subject "SERVER" --unit "C"
```

This returns the metric object

```
{
 "id": 676,
 "version": 1,
 "state": "ENABLED",
 "name": "Core temperature",
 "description": "",
 "unit": "C",
 "period": 60,
 "source": "cli",
 "dataType": "DOUBLE",
 "subject": "SERVER"
}
```

Push data for the server metric. In this example we will push data for the 'Core temperature' metric (id 676) for the 'myserver' server (id 34) at timestamp 1495108650 with value 50.4:

```
coscale-cli data insert --data="M676:S34:1495108650:50.4"
```


[For more information, check out the CLI documentation.](http://docs.coscale.com/tools/cli/index/)