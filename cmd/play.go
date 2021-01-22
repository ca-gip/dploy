/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	ansibler "github.com/apenella/go-ansible"
	"github.com/ca-gip/dploy/internal/ansible"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Run ansible-playbook command",
	Long:  `TODO...`,
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := os.Getwd()
		play(cmd, args, path)
	},
}

func init() {
	rootCmd.AddCommand(playCmd)

	// Required flags
	playCmd.Flags().StringSliceP("filter", "", nil, `filters inventory based its on vars ex: "foo==bar,bar!=foo""`)
	_ = playCmd.MarkFlagRequired("filter")
	playCmd.Flags().StringP("playbook", "p", "", "playbook to run")
	_ = playCmd.MarkFlagRequired("playbook")

	// Ansible params
	playCmd.Flags().StringP("vault-password-file", "", "", "vault password file")
	playCmd.Flags().StringSliceP("limit", "l", nil, "further limit selected hosts to an additional pattern")
	playCmd.Flags().StringSliceP("tags", "t", nil, "only run plays and tasks tagged with these values")
}

func play(cmd *cobra.Command, args []string, path string) {
	// Load project from root
	project := ansible.Projects.LoadFromPath(path)

	// Retrieve playbook to be run
	playbookPath, _ := cmd.Flags().GetString("playbook")
	playbook, err := project.PlaybookPath(playbookPath)
	if err != nil {
		log.Fatalf(`%s not a valid path`, playbookPath)
	}

	// Retrieve filter to select inventories
	rawFilters, _ := cmd.Flags().GetStringSlice("filter")
	filters := ansible.ParseFilterArgsFromSlice(rawFilters)
	inventories := project.FilterInventory(filters)

	// Retrieve ansible flags
	tags, _ := cmd.Flags().GetStringSlice("tags")
	limit, _ := cmd.Flags().GetStringSlice("limit")
	vaultPassFile, _ := cmd.Flags().GetString("vault-password-file")

	// Execute ansible for each inventory (sequential)
	for _, inventory := range inventories {
		ansiblePlaybookOptions := &ansibler.AnsiblePlaybookOptions{
			Inventory:         inventory.RelativePath(),
			Limit:             strings.Join(limit, ","),
			Tags:              strings.Join(tags, ","),
			VaultPasswordFile: vaultPassFile,
		}
		play := ansibler.AnsiblePlaybookCmd{
			Playbook: playbook.RelativePath(),
			Options:  ansiblePlaybookOptions,
		}

		ansibler.AnsibleForceColor()
		err := play.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
