package command

import (
	"coscale/api"
	"os"
)

// CheckObjectName is the name of the check-config subcommand
var CheckObjectName = "check-config"

// CheckObject is the check-config subcommand and is used to check to api configuration
var CheckObject = &Command{
	Name:      CheckObjectName,
	UsageLine: `check-config is used to check to api configuration file`,
	Run: func(cmd *Command, args []string) {
		// check for getting the config file path
		file, err := GetConfigPath()
		if err != nil {
			cmd.PrintResult("", err)
		}
		// check if the file actually exists
		if _, err := os.Stat(file); err != nil {
			cmd.PrintResult("", err)
		}
		// check if the configuration file can be parsed
		config, err := api.ReadApiConfiguration(file)
		if err != nil {
			cmd.PrintResult("", err)
		}
		// check if we can loggin with this configuration
		api := api.NewApi(config.BaseUrl, config.AccessToken, config.AppId, false)
		err = api.Login()
		if err != nil {
			cmd.PrintResult("", err)
		}
		cmd.PrintResult(`{"msg":"Configuration successfuly checked!"}`, nil)
	},
}
