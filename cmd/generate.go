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
	"github.com/ca-gip/dploy/internal/services"
	"github.com/spf13/cobra"
	"os"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate ansible-playbook command for multiple inventories",
	Long: `
TODO`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(cmd.Flags().GetStringArray("select"))

		fmt.Println("generate called")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringSliceP("filter", "", nil, "select inventory based on vars ex: foo==bar,bar!=foo")

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

	home, _ := os.UserHomeDir()
	path := fmt.Sprintf("%s/%s", home, "Projects/ansible-kube/inventories")
	k8s := services.LoadFromPath(path)
	_ = generateCmd.RegisterFlagCompletionFunc("filter", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return k8s.GetInventoryKeys(), cobra.ShellCompDirectiveDefault
	})

}
