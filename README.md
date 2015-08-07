# CoScale command line interface
The CoScale command line interface is an interface to the CoScale API.

It provides easy to use methods for:
* Managing servers and server groups.
* Push metric groups, metrics and data.
* Register events and event categories.
* List, acknowledge and resolve your alerts.

## Compilation

### Linux / Mac OS X
![CoScale CLI installation](install.gif)

1. `$ git clone https://github.com/CoScale/coscale-tools.git`
1. [(Install Go Programming Language)](https://golang.org/doc/install)
1. `$ cd coscale-tools/cli`
1. `sh build.sh`

## Usage
    cli-cmd a tool for CoScale Api.

    Usage:
    alert <action> [--<field>='<data>']

    Actions for command "alert":

    list
            alert list [--filter]
    acknowledge
            alert acknowledge (--id | --name)
    resolve
            alert resolve (--id | --name)


    The json objects are returned formatted by default, but can be returned on 1 line by using:
    --rawOutput

    The CoScale api configuration (authentication) by default will be taken from api.conf file,
    placed in the same folder with the cli-cmd. api.conf file it is the same configuration file
    used by the CoScale agent. If the api.conf file doesn't exists, the informations also can be
    provided on the command line using:
    --api-url
        Base url for the api (optional, default = "https://api.coscale.com/").
    --app-id
        The application id.
    --access-token
        A valid access token for the given application.

    Use "cli-cmd [object] <help>" for more information about a command.

## Examples

Adding a CoScale event category with custom attributes exitCode and executionTime (in seconds).

`./coscale-cli event new --name "Example category" --attributeDescriptions "[{\"name\":\"exitCode\", \"type\":\"integer\"}, {\"name\":\"executionTime\", \"type\":\"integer\", \"unit\":\"s\"}]" --source "CLI"`

Adding a CoScale event to "Example category" with name 'Event example', exitCode '0' and an executionTime of 10 seconds.

`./coscale-cli event data --name "Example category" --message "Event example" --subject "a" --attribute "{\"exitCode\":0, \"executionTime\":10}"``
