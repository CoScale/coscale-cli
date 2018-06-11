package command

import (
	"coscale/api"
	"flag"
	"fmt"
	"os"
	"strings"
)

// configObjectName is the name of the config subcommand
var configObjectName = "config"

// ConfigObject defines the config command on the CLI.
var ConfigObject = NewCommand(configObjectName, "config <action> [--<field>='<data>']", ConfigActions)

// ConfigActions defines the config actions on the CLI.
var ConfigActions = []*Command{
	{
		Name:      "check",
		UsageLine: "config check",
		Long: `
Check the CLI configuration.
`,
		Run: func(cmd *Command, args []string) {
			// check for getting the config file path
			file, err := GetConfigPath()
			if err != nil {
				fmt.Fprintln(os.Stderr, "No such file: "+file)
				os.Exit(EXIT_SUCCESS_ERROR)
			}
			// check if the configuration file can be parsed
			config, err := api.ReadApiConfiguration(file)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Could not parse configuration from "+file)
				os.Exit(EXIT_SUCCESS_ERROR)
			}
			// check if we can loggin with this configuration
			api := api.NewApi(config.BaseUrl, config.AccessToken, config.AppId, false, false)
			err = api.Login()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Api authentication failed")
				os.Exit(EXIT_SUCCESS_ERROR)
			}
			fmt.Fprintln(os.Stderr, "Configuration successfully checked")
			os.Exit(EXIT_SUCCESS)
		},
	},
	{
		Name:      "set",
		UsageLine: "config set (--api-url --app-id --access-token)",
		Long: `
Write the CLI configuration file, a file api.conf will be created in the same directory as
the coscale-cli binary.

Mandatory:
	--api-url
		Base url for the api (optional, default = "https://api.coscale.com").
	--app-id
		The application id.
	--access-token
		A valid access token for the given application.
`,
		Run: func(cmd *Command, args []string) {
			// get the config path
			path, _ := GetConfigPath()
			if path == "" {
				fmt.Fprintln(os.Stderr, "Could not determine the CLI configuration path.")
				os.Exit(EXIT_SUCCESS_ERROR)
			}
			// create the config json
			var baseUrl, accessToken, appId string

			var flags flag.FlagSet
			flags.StringVar(&baseUrl, "api-url", "https://api.coscale.com", "Base url for the api")
			flags.StringVar(&appId, "app-id", "", "The application id")
			flags.StringVar(&accessToken, "access-token", "", "A valid access token for the given application")
			flags.Parse(args)

			config := &api.ApiConfiguration{strings.Trim(baseUrl, "/"), accessToken, appId}

			if config.BaseUrl == "" || config.AccessToken == "" || config.AppId == "" {
				cmd.PrintUsage()
				os.Exit(EXIT_FLAG_ERROR)
			}

			err := api.WriteApiConfiguration(path, config)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to write the configuration file.")
				os.Exit(EXIT_SUCCESS_ERROR)
			}
			// write the json to the file
			fmt.Fprintln(os.Stderr, "Successfully wrote CLI configuration file.")
			os.Exit(EXIT_SUCCESS)
		},
	},
}
