// Copyright (c) 2016 ECS Team, Inc. - All Rights Reserved
// https://github.com/ECSTeam/cloudfoundry-top-plugin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/cloudfoundry/cli/plugin"
	"os"
	"regexp"
	"strings"
	"github.com/simonleung8/flags"
	"bytes"
	"encoding/json"
	"time"
)


// Events represents Buildpack Usage CLI interface
type Events struct{}

// OutputResults represents the filtered event results for the input args
type OutputResults struct {
	Comment      string  `json:"comment"`
	Resources    []AppSearchResources `json:"resources"`
}

// Metadata is the data retrived from the response json
type Metadata struct {
	GUID string `json:"guid"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Inputs represent the parsed input args
type Inputs struct {
	fromDate time.Time
	toDate   time.Time
	isCsv    bool
	isJson   bool
}

// GetMetadata provides the Cloud Foundry CLI with metadata to provide user about how to use `inspect` command
func (c *Events) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "inspect",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 1,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "inspect",
				HelpText: "Interrogate and inspect CF foundation (by akoranne@ecsteam.com)",
				UsageDetails: plugin.Usage {
					Usage: UsageText(),
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(Events))
}

// Run is what is executed by the Cloud Foundry CLI when the inspect command is specified
func (c Events) Run(cli plugin.CliConnection, args []string) {
	var ins Inputs

	switch args[0] {
	case "inspect":
		ins = c.buildClientOptions(args)
	case "example-alternate-command":
	default:
		return
	}

	orgs := c.GetOrgs(cli)
	spaces := c.GetSpaces(cli)
	apps := c.GetAppData(cli)
	results := c.FilterResults(cli, ins, orgs, spaces, apps)
	if (ins.isCsv) {
		c.AppsInCSVFormat(results)
	} else {
		c.AppsInJsonFormat(results)
	}
}


func Usage(code int) {
	fmt.Println("\nUsage: ", UsageText())
	os.Exit(code)
}

func UsageText() (string) {
	usage := "cf inspect [--json]" +
		"\n       --json : list output in json format (default is csv)\n"
	return usage
}

// sanitize data by replacing \r, and \n with ';'
func sanitize(data string) (string) {
	var re = regexp.MustCompile(`\r?\n`)
	var str = re.ReplaceAllString(data, ";")
	str = strings.Replace(str, ";;", ";", 1)
	return str;
}

// read arguments passed for the plugin
func (c *Events) buildClientOptions(args[] string) (Inputs) {
	fc := flags.New()
	fc.NewBoolFlag("json", "js", "list output in json format (default is csv)")

	//err := fc.Parse(args[1:]...)
	//
	//if err != nil {
	//	fmt.Println("\n Receive error reading arguments ... ", err)
	//	Usage(1)
	//}

	var ins Inputs
	ins.isCsv = true
	ins.isJson = false

	if (fc.IsSet("json")) {
		ins.isJson = true
		ins.isCsv = false
	}

	return ins
}

// prints the results as a csv text to console
func (c Events) AppsInCSVFormat(results OutputResults) {
	fmt.Println("")
	fmt.Printf(results.Comment)
	fmt.Printf("%s,%s,%s,%s,%s,%s\n", "AppName", "AppGUID", "OrgName", "OrgGUID", "SpaceName", "SpaceGUID")
	for _, val := range results.Resources  {
		fmt.Printf("%s,%s,%s,%s,%s,%s\n",
			val.Entity.Name, val.Metadata.GUID,
			val.Entity.Org, val.Entity.OrgGUID,
			val.Entity.Space, val.Entity.SpaceGUID)
	}
}

// prints the results as a json text to console
func (c Events) AppsInJsonFormat(results OutputResults) {
	var out bytes.Buffer
	b, _ := json.Marshal(results)
	err := json.Indent(&out, b, "", "\t")
	if err != nil {
		fmt.Println(" Recevied error formatting json output.")
	} else {
		fmt.Println(out.String())
	}
}

