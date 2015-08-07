package command

import (
	"coscale/api"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

const (
	DEFAULT_FLAG_VALUE string = `!>dUmmy<!`
)

const (
	// Status codes for os.Exit.
	EXIT_SUCCESS              int = 0
	EXIT_SUCCESS_ERROR        int = 1
	EXIT_AUTHENTICATION_ERROR int = 2
	EXIT_FLAG_ERROR           int = 3
)

type Command struct {
	Name        string
	UsageLine   string
	Long        string
	SubCommands []*Command
	Capi        *api.Api //api connector
	Flag        flag.FlagSet
	Run         func(cmd *Command, args []string)
}

func NewCommand(name, usage string, subCommands []*Command) *Command {
	return &Command{
		Name:        name,
		UsageLine:   usage,
		SubCommands: subCommands,
		Run: func(cmd *Command, args []string) {
			subCmd := cmd.GetSubCommand(args)
			if subCmd != nil {
				subCmd.Run(subCmd, args[1:])
			}
		},
	}
}

func (c *Command) Runnable() bool {
	return len(c.SubCommands) == 0
}

func (c *Command) GetSubCommand(args []string) *Command {
	if len(args) < 1 {
		c.PrintUsage()
	}
	for _, cmd := range c.SubCommands {
		if cmd.Name == args[0] {
			return cmd
		}
	}
	c.PrintUsage()
	return nil
}

func (c *Command) GetAllSubCommands() []*Command {
	commands := make([]*Command, 0, 0)
	if c.Runnable() {
		commands = append(commands, c)
	} else {
		for _, subCmd := range c.SubCommands {
			commands = append(commands, subCmd.GetAllSubCommands()...)
		}
	}
	return commands
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

//make the first letter form a string uppercase
func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func (c *Command) PrintUsage() {
	tmpl(os.Stderr, usageTemplate, c)
	os.Exit(2)
}

//return a Api object
func (c *Command) GetApi(baseUrl, accessToken, appId string, rawOutput bool) *api.Api {
	if accessToken == "" || appId == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			os.Exit(EXIT_FLAG_ERROR)
		}
		config, err := api.ReadApiConfiguration(dir + "/api.conf")
		if err != nil {
			c.PrintUsage()
		}
		baseUrl = config.BaseUrl
		accessToken = config.AccessToken
		appId = config.AppId
	}
	return api.NewApi(baseUrl, accessToken, appId, rawOutput)
}

func (c *Command) ParseArgs(args []string) {
	if len(args) > 0 && args[0] == "help" {
		c.PrintUsage()
	}
	//add the flags for the api configuration
	var baseUrl, accessToken, appId string
	var rawOutput bool
	c.Flag.StringVar(&baseUrl, "api-url", "https://api.coscale.com", "Base url for the api")
	c.Flag.StringVar(&appId, "app-id", "", "The application id")
	c.Flag.StringVar(&accessToken, "access-token", "", "A valid access token for the given application")
	c.Flag.BoolVar(&rawOutput, "rawOutput", false, "The returned json objects are returned formatted by default")
	c.Flag.Parse(args)
	unknownArgs := c.Flag.Args()
	if len(unknownArgs) > 0 && unknownArgs[0] != "help" {
		fmt.Fprintf(os.Stderr, "Unknown field %s\n", unknownArgs[0])
		os.Exit(EXIT_FLAG_ERROR)
	}
	c.Capi = c.GetApi(baseUrl, accessToken, appId, rawOutput)
}

func (c *Command) PrintResult(result string, err error) {
	if err == nil {
		fmt.Fprint(os.Stdout, result)
		os.Exit(EXIT_SUCCESS)
	} else if api.IsAuthenticationError(err) {
		fmt.Fprint(os.Stderr, `{"msg":"Authentication failed!"}`)
		os.Exit(EXIT_AUTHENTICATION_ERROR)
	} else {
		fmt.Fprintf(os.Stderr, GetErrorJson(err))
		os.Exit(EXIT_SUCCESS_ERROR)
	}
}

// GetErrorJson return only the json string from a error message from api
func GetErrorJson(err error) string {
	if strings.Index(err.Error(), `{`) > -1 {
		return err.Error()[strings.Index(err.Error(), `{`):]
	}
	return err.Error()
}

var usageTemplate = `cli-cmd a tool for CoScale Api.

Usage:
	{{.UsageLine}}
{{if .Runnable}}
{{.Name | printf "Action \"%s\" usage:"}} 

{{.Long | trim}}{{else}}
{{.Name | printf "Actions for command \"%s\":"}}
{{range .SubCommands}}
	{{.Name | printf "%s"}}
			{{.UsageLine | printf "%-11s"}}{{end}}
    {{end}}

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
`
