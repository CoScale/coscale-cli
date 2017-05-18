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

Don't forget to fill in your *application id* and *access token* as provided on the Access Token page in the CoScale UI.

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

1) Create a basic CoScale event category called “Releases”

```
coscale-cli event new --name "Releases"
```

2) Add a CoScale event about to the “Releases” event category

```
coscale-cli event newdata --name "Releases" --message "V2.3.4" --subject "a"
```


### Server Examples

1. Add a server that will not run a CoScale agent, but can be used to attach events to.

```
coscale-cli server new --name "cron server" --description "Server that that runs cron jobs"
```

2. Delete the cron server.

```
coscale-cli server delete --name "cron server"
```


### Alert Examples

1. View all unresolved alerts.

```
coscale-cli alert list --filter unresolved
```

2. Resolve an alert from the list. The id is the id of an alert we got from the list command.


```
coscale-cli alert resolve --id 347482
```


[For more information, check out the CLI documentation.](http://docs.coscale.com/tools/cli/index/)
