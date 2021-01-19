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
	"fmt"
	"github.com/ca-gip/dploy/internal/ansible"
	"github.com/spf13/cobra"
	"log"
	"os"
	"regexp"
	"strings"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate ansible-playbook command",
	Long: `
TODO`,
	Run: func(cmd *cobra.Command, args []string) {

		curr, _ := os.Getwd()
		project := ansible.LoadFromPath(curr)

		rawFilters, _ := cmd.Flags().GetStringSlice("filter")
		filters := ansible.ParseFilterArgsFromSlice(rawFilters)
		inventories := project.Inventories.Filter(filters)

		playbookPath, _ := cmd.Flags().GetString("playbook")
		playbook := project.Playbooks.GetPlaybook(playbookPath)

		if playbook == nil {
			log.Fatalf(`%s not a valid path`, playbookPath)
		}

		tags, _ := cmd.Flags().GetStringSlice("tags")
		limit, _ := cmd.Flags().GetStringSlice("limit")
		skipTags, _ := cmd.Flags().GetStringSlice("skip-tags")
		check, _ := cmd.Flags().GetBool("check")
		diff, _ := cmd.Flags().GetBool("diff")
		vaultPassFile, _ := cmd.Flags().GetString("vault-password-file")
		askVaultPass, _ := cmd.Flags().GetBool("ask-vault-password")

		commands := &ansible.AnsibleCommandTpl{
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

	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringSliceP("filter", "", nil, `filters inventory based its on vars ex: "foo==bar,bar!=foo""`)
	generateCmd.Flags().StringP("playbook", "p", "", "playbook to run")
	_ = generateCmd.MarkFlagRequired("filter")
	_ = generateCmd.MarkFlagRequired("playbook")

	// Ansible params
	generateCmd.Flags().BoolP("ask-vault-password", "", false, "ask for vault password")
	generateCmd.Flags().StringP("vault-password-file", "", "", "vault password file")
	generateCmd.Flags().StringSliceP("skip-tags", "", nil, "only run plays and tasks whose tags do not match these values")
	generateCmd.Flags().BoolP("check", "C", false, "don't make any changes; instead, try to predict some of the changes that may occur")
	generateCmd.Flags().BoolP("diff", "D", false, "when changing (small) files and templates, show the differences in those files; works great with --check")
	generateCmd.Flags().StringSliceP("limit", "l", nil, "further limit selected hosts to an additional pattern")
	generateCmd.Flags().StringSliceP("tags", "t", nil, "only run plays and tasks tagged with these values")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	_ = generateCmd.RegisterFlagCompletionFunc("filter", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		key, op, value := ansible.ParseFilter(toComplete)

		cobra.CompDebug(fmt.Sprintf("key:%s op:%s value:%s", key, op, value), true)

		curr, _ := os.Getwd()
		k8s := ansible.LoadFromPath(curr)

		availableKeys := k8s.Inventories.GetInventoryKeys()

		blank := key == "" && op == "" && value == ""
		if blank {
			return availableKeys, cobra.ShellCompDirectiveDefault
		}

		writingKey := key != "" && op == "" && value == ""
		if writingKey {
			var keysCompletion []string
			for _, availableKey := range availableKeys {
				if strings.HasPrefix(availableKey, key) {
					keysCompletion = append(keysCompletion, availableKey)
				}
			}

			if len(keysCompletion) == 1 {
				var prefixedOperator []string

				for _, allowedOperator := range ansible.AllowedOperators {
					prefixedOperator = append(prefixedOperator, fmt.Sprintf("%s%s", keysCompletion[0], allowedOperator))
				}
				return prefixedOperator, cobra.ShellCompDirectiveDefault
			}

			return keysCompletion, cobra.ShellCompDirectiveDefault
		}

		writingOp := key != "" && op != "" && value == ""
		if writingOp {
			var prefixedOperator []string

			for _, allowedOperator := range ansible.AllowedOperators {

				if op == allowedOperator {
					availableValues := k8s.Inventories.GetInventoryValues(key)

					var prefixedValues []string

					for _, availableValue := range availableValues {

						if availableValue != "" {
							prefixedValues = append(prefixedValues, fmt.Sprintf("%s%s%s", key, op, availableValue))
						}

					}

					return prefixedValues, cobra.ShellCompDirectiveDefault

				}

				if allowedOperator[0] == op[0] {
					prefixedOperator = append(prefixedOperator, fmt.Sprintf("%s%s", key, allowedOperator))
				}

			}

			if len(prefixedOperator) == 1 {
				availableValues := k8s.Inventories.GetInventoryValues(key)

				_, foundOp, _ := ansible.ParseFilter(prefixedOperator[0])

				var prefixedValues []string

				for _, availableValue := range availableValues {

					if availableValue != "" {
						prefixedValues = append(prefixedValues, fmt.Sprintf("%s%s%s", key, foundOp, availableValue))
					}

				}

				return prefixedValues, cobra.ShellCompDirectiveDefault
			}

			return prefixedOperator, cobra.ShellCompDirectiveDefault
		}

		writingValue := key != "" && op != "" && value != ""
		if writingValue {
			for _, allowedOperator := range ansible.AllowedOperators {

				if op == allowedOperator {
					availableValues := k8s.Inventories.GetInventoryValues(key)

					var prefixedValues []string

					for _, availableValue := range availableValues {
						if availableValue != "" && strings.HasPrefix(availableValue, value) {
							prefixedValues = append(prefixedValues, fmt.Sprintf("%s%s%s", key, op, availableValue))
						}

					}

					return prefixedValues, cobra.ShellCompDirectiveDefault

				}

			}
			return []string{}, cobra.ShellCompDirectiveDefault

		}

		return k8s.Inventories.GetInventoryKeys(), cobra.ShellCompDirectiveDefault

	})

	_ = generateCmd.RegisterFlagCompletionFunc("playbook", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		curr, _ := os.Getwd()
		k8s := ansible.LoadFromPath(curr)
		return k8s.Playbooks.GetPlaybooks(), cobra.ShellCompDirectiveDefault
	})

	_ = generateCmd.RegisterFlagCompletionFunc("tags", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		var _ = regexp.MustCompile("([\\w-.\\/]+)([,]|)")

		playbookPath, _ := cmd.Flags().GetString("playbook")
		if len(playbookPath) == 0 {
			return nil, cobra.ShellCompDirectiveDefault
		}

		curr, _ := os.Getwd()
		project := ansible.LoadFromPath(curr)
		playbook := project.Playbooks.GetPlaybook(playbookPath)

		return playbook.AllTags.List(), cobra.ShellCompDirectiveDefault

	})

}
