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
	"context"
	"github.com/apenella/go-ansible/pkg/adhoc"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/ca-gip/dploy/internal/ansible"
	"github.com/ca-gip/dploy/internal/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run Ad Hoc command",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		curr, _ := os.Getwd()
		exec(cmd, args, curr)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	// Required flags
	execCmd.Flags().StringSliceP("filter", utils.EmptyString, nil, `filters inventory based its on vars ex: "foo==bar,bar!=""`)
	_ = execCmd.MarkFlagRequired("filter")

	// Ansible params
	execCmd.Flags().StringP("module", "m", utils.EmptyString, "module name to execute (default=command)")
	execCmd.Flags().StringP("pattern", "p", utils.EmptyString, "host pattern")
	execCmd.Flags().StringP("args", "a", utils.EmptyString, "module arguments")
	_ = execCmd.MarkFlagRequired("args")
	execCmd.Flags().StringSliceP("extra-vars", "e", []string{}, "set additional variables as key=value or YAML/JSON, if filename prepend with @")

	execCmd.Flags().IntP("background", "B", 0, "run asynchronously, failing after X seconds")

	// Completions
	_ = execCmd.RegisterFlagCompletionFunc("filter", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		path, _ := os.Getwd()
		return filterCompletion(toComplete, path)
	})

	_ = execCmd.RegisterFlagCompletionFunc("pattern", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		path, _ := os.Getwd()
		return hostPatternCompletion(toComplete, path)
	})

}

func template(cmd *cobra.Command, args []string, path string) {
	// Load project from root
	project := ansible.Projects.LoadFromPath(path)

	// Retrieve filter to select inventories
	rawFilters, _ := cmd.Flags().GetStringSlice("filter")
	filters := ansible.ParseFilterArgsFromSlice(rawFilters)
	inventories := project.FilterInventory(filters)

	// Retrieve ansible flags
	module, _ := cmd.Flags().GetString("module")
	pattern, _ := cmd.Flags().GetString("pattern")
	arg, _ := cmd.Flags().GetString("args")
	extra, _ := cmd.Flags().GetStringSlice("extra-vars")
	background, _ := cmd.Flags().GetInt("background")

	//
	vars, err := varListToMap(extra)
	if err != nil {
		log.Fatal(err)
	}

	template := ansible.AdHocCmd{
		Comment:           "# Commands",
		Inventory:         inventories,
		Pattern:           pattern,
		ModuleName:        module,
		ModuleArgs:        arg,
		ExtraVars:         vars,
		Background:        background,
		Fork:              0,
		PollInterval:      0,
		Limit:             nil,
		Check:             false,
		Diff:              false,
		OneLine:           false,
		Tree:              false,
		PlaybookDir:       "",
		VaultPasswordFile: "",
		AskVaultPass:      false,
	}

	template.Generate()
}

func exec(cmd *cobra.Command, args []string, path string) {
	// Load project from root
	project := ansible.Projects.LoadFromPath(path)

	// Retrieve filter to select inventories
	rawFilters, _ := cmd.Flags().GetStringSlice("filter")
	filters := ansible.ParseFilterArgsFromSlice(rawFilters)
	inventories := project.FilterInventory(filters)

	// Retrieve ansible flags
	module, _ := cmd.Flags().GetString("module")
	pattern, _ := cmd.Flags().GetString("pattern")
	arg, _ := cmd.Flags().GetString("args")
	background, _ := cmd.Flags().GetInt("background")
	extra, _ := cmd.Flags().GetStringSlice("extra-vars")

	vars, err := varListToMap(extra)
	if err != nil {
		log.Fatal(err)
	}

	for _, inventory := range inventories {
		ansibleAdhocOptions := &adhoc.AnsibleAdhocOptions{
			Args:       arg,
			Background: background,
			ExtraVars:  vars,
			Inventory:  inventory.RelativePath(),
			ModuleName: module,
		}

		adhoc := &adhoc.AnsibleAdhocCmd{
			Pattern: pattern,
			Options: ansibleAdhocOptions,
		}

		options.AnsibleForceColor()
		err := adhoc.Run(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}

}
