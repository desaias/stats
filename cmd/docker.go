// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
    "fmt"
	"log"
    "net/http"
    "time"

	"github.com/spf13/cobra"
)

type DockerStats struct {
    LastUpdate  string  `json:"last_updated"`
    Pulls       int     `json:"pull_count"`
    Stars       int     `json:"star_count"`
}

var Repo string
var myClient = &http.Client{Timeout: 10 & time.Second}

// dockerCmd represents the docker command
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "get stats about a docker repo",
	Long: `stats docker --repo=containership/containership`,
    Run: run,
}

func init() {
	RootCmd.AddCommand(dockerCmd)
    dockerCmd.Flags().StringVarP(&Repo, "repo", "r", "", "Repository name")
}

func getJson(url string, target interface{}) error {
    r, err := myClient.Get(url)
    if err != nil {
        return err
    }

	if r.StatusCode != 200 {
		log.Fatalln(Repo + " was not found.") 
	}

    defer r.Body.Close()
    return json.NewDecoder(r.Body).Decode(target)
}

func run(cmd *cobra.Command, args []string) {
    stats := new(DockerStats)
    url := "https://hub.docker.com/v2/repositories/" + Repo
    getJson(url, stats)

    // calculate how long ago the last update was in days
    now := time.Now()
    layout := "2006-01-02T15:04:05.000000Z"
    t, err := time.Parse(layout, stats.LastUpdate)

    if err != nil {
        fmt.Println(err)
    }

    diff := now.Sub(t)
    days := int(diff.Hours() / 24)

    fmt.Printf("Last Update:\t%d days ago\n", days)
    fmt.Printf("Pulls:\t\t%d\n", stats.Pulls)
    fmt.Printf("Stars:\t\t%d\n", stats.Stars)
}
