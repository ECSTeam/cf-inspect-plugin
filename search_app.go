package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

// AppSearchResults represents top level attributes of JSON response from Cloud Foundry API
type AppSearchResults struct {
	TotalResults int                  `json:"total_results"`
	TotalPages   int                  `json:"total_pages"`
	PrevUrl      string               `json:"prev_url"`
	NextUrl      string               `json:"next_url"`
	Resources    []AppSearchResources `json:"resources"`
}

// AppSearchResources represents resources attribute of JSON response from Cloud Foundry API
type AppSearchResources struct {
	Entity   AppSearchEntity `json:"entity"`
	Metadata Metadata        `json:"metadata"`
}

// AppSearchEntity represents entity attribute of resources attribute within JSON response from Cloud Foundry API
type AppSearchEntity struct {
	Name              string `json:"name"`
	Buildpack         string `json:"buildpack"`
	DetectedBuildpack string `json:"detected_buildpack"`
	Instances         int    `json:"instances"`
	State             string `json:"state"`
	Memory            int    `json:"memory"`
	DiskQuota         int    `json:"disk_quota"`
	SpaceGUID         string `json:"space_guid"`
	OrgGUID           string `json:"organization_guid"`
	Space             string `json:"space"`
	Org               string `json:"org"`
}

//// GetAppData requests all of the Application data from Cloud Foundry
//func (c Events) GetApps(cli plugin.CliConnection) map[string]AppSearchResults {
//	var data map[string]AppSearchResults
//	data = make(map[string]AppSearchResults)
//	spaces := c.GetAppData(cli)
//
//	for _, val := range spaces.Resources {
//		data[val.Metadata.GUID] = val.Metadata
//	}
//
//	return data
//}

// GetAppData requests all of the Application data from Cloud Foundry
func (c Events) GetAppData(cli plugin.CliConnection) AppSearchResults {
	var res AppSearchResults
	res = c.UnmarshallAppSearchResults("/v2/apps?order-direction=asc&results-per-page=100", cli)

	if res.TotalPages > 1 {
		for i := 2; i <= res.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v2/apps?order-direction=asc&page=%v&results-per-page=100", strconv.Itoa(i))
			tRes := c.UnmarshallAppSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func (c Events) UnmarshallAppSearchResults(apiUrl string, cli plugin.CliConnection) AppSearchResults {
	var tRes AppSearchResults
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}

// filter the results for given input criteria.
func (c Events) FilterResults(cli plugin.CliConnection, ins Inputs, orgs map[string]string, spaces map[string]SpaceSearchEntity, apps AppSearchResults) (OutputResults) {
	var results OutputResults

	type AppSearchEntity struct {
		Name              string `json:"name"`
		Buildpack         string `json:"buildpack"`
		DetectedBuildpack string `json:"detected_buildpack"`
		Instances         int    `json:"instances"`
		State             string `json:"state"`
		Memory            int    `json:"memory"`
		DiskQuota         int    `json:"disk_quota"`
		SpaceGUID         string `json:"space_guid"`
		OrgGUID           string `json:"organization_guid"`
		Space             string `json:"space"`
		Org               string `json:"org"`
	}


	for _, val := range apps.Resources {
		var outEntity AppSearchResources

		outEntity.Metadata  = val.Metadata
		outEntity.Entity.Name = val.Entity.Name
		outEntity.Entity.Buildpack         = val.Entity.Buildpack
		outEntity.Entity.DetectedBuildpack = val.Entity.DetectedBuildpack
		outEntity.Entity.Instances         = val.Entity.Instances
		outEntity.Entity.State             = val.Entity.State
		outEntity.Entity.Memory            = val.Entity.Memory
		outEntity.Entity.DiskQuota         = val.Entity.DiskQuota
		outEntity.Entity.SpaceGUID         = val.Entity.SpaceGUID
		outEntity.Entity.OrgGUID           = spaces[val.Entity.SpaceGUID].OrgGUID
		outEntity.Entity.Org               = orgs[spaces[val.Entity.SpaceGUID].OrgGUID]
		outEntity.Entity.Space             = spaces[val.Entity.SpaceGUID].Name
		results.Resources = append(results.Resources, outEntity)

	}
	return results;
}
