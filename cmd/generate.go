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
	"github.com/ca-gip/dploy/internal/ansible"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate ansible-playbook command",
	Long:  `TODO...`,
	Run: func(cmd *cobra.Command, args []string) {
		curr, _ := os.Getwd()
		generate(cmd, args, curr)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Required flags
	generateCmd.Flags().StringSliceP("filter", "", nil, `filters inventory based its on vars ex: "foo==bar,bar!=foo""`)
	_ = generateCmd.MarkFlagRequired("filter")
	generateCmd.Flags().StringP("playbook", "p", "", "playbook to run")
	_ = generateCmd.MarkFlagRequired("playbook")

	// Ansible params
	generateCmd.Flags().BoolP("ask-vault-password", "", false, "ask for vault password")
	generateCmd.Flags().StringP("vault-password-file", "", "", "vault password file")
	generateCmd.Flags().StringSliceP("skip-tags", "", nil, "only run plays and tasks whose tags do not match these values")
	generateCmd.Flags().BoolP("check", "C", false, "don't make any changes; instead, try to predict some of the changes that may occur")
	generateCmd.Flags().BoolP("diff", "D", false, "when changing (small) files and templates, show the differences in those files; works great with --check")
	generateCmd.Flags().StringSliceP("limit", "l", nil, "further limit selected hosts to an additional pattern")
	generateCmd.Flags().StringSliceP("tags", "t", nil, "only run plays and tasks tagged with these values")

	// Completions
	_ = generateCmd.RegisterFlagCompletionFunc("filter", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		path, _ := os.Getwd()
		return filterCompletion(toComplete, path)
	})

	_ = generateCmd.RegisterFlagCompletionFunc("playbook", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		path, _ := os.Getwd()
		return playbookCompletion(toComplete, path)
	})

	_ = generateCmd.RegisterFlagCompletionFunc("tags", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		path, _ := os.Getwd()
		playbookPath, _ := cmd.Flags().GetString("playbook")
		return tagsCompletion(toComplete, path, playbookPath)
	})

}

func generate(cmd *cobra.Command, args []string, path string) {
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
	skipTags, _ := cmd.Flags().GetStringSlice("skip-tags")
	check, _ := cmd.Flags().GetBool("check")
	diff, _ := cmd.Flags().GetBool("diff")
	vaultPassFile, _ := cmd.Flags().GetString("vault-password-file")
	askVaultPass, _ := cmd.Flags().GetBool("ask-vault-password")

	commands := &ansible.Command{
		Inventory:         inventories,
		Playbook:          playbook,
		Tags:              tags,
		Limit:             limit,
		SkipTags:          skipTags,
		Check:             check,
		Diff:              diff,
		VaultPasswordFile: vaultPassFile,
		AskVaultPass:      askVaultPass,
	}

	commands.GenerateCmd()
}
