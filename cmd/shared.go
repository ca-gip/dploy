package cmd

import (
	"fmt"
	"github.com/ca-gip/dploy/internal/ansible"
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

func extractMultipleCompletion(toComplete string) (remainder string, current string) {
	toCompletes := strings.Split(toComplete, ",")

	if len(toCompletes) == 1 {
		return utils.EmptyString, toCompletes[0]
	}

	remainder = fmt.Sprintf("%s%s", strings.Join(toCompletes[:len(toCompletes)-1], ","), ",")
	current = toCompletes[len(toCompletes)-1]
	return
}

func filterCompletion(toComplete string, path string) ([]string, cobra.ShellCompDirective) {
	logrus.SetLevel(logrus.PanicLevel)
	var EmptyCompletion []string

	remainder, current := extractMultipleCompletion(toComplete)
	cobra.CompDebug(fmt.Sprintf("extract muitple: remainder:%s current:%s\n", remainder, current), true)

	key, op, value := ansible.ParseFilter(current)
	cobra.CompDebug(fmt.Sprintf("parse filter: key:%s op:%s value:%s\n", key, op, value), true)

	k8s := ansible.Projects.LoadFromPath(path)
	availableKeys := k8s.InventoryKeys()

	// Blank
	blank := key == utils.EmptyString && op == utils.EmptyString && value == utils.EmptyString
	if blank {
		cobra.CompDebug("blank case", true)
		return utils.AppendPrefixOnSlice(remainder, availableKeys), cobra.ShellCompDirectiveDefault
	}

	// Writing Key
	writingKey := key != utils.EmptyString && op == utils.EmptyString && value == utils.EmptyString
	if writingKey {
		cobra.CompDebug("writingKey case", true)
		keyMatches := utils.Filter(availableKeys, func(s string) bool {
			return strings.HasPrefix(s, key)
		})
		// Key is known complete for operators
		if len(keyMatches) == 1 {
			key := keyMatches[0]
			keyOperatorCompletion := utils.AppendPrefixOnSlice(key, ansible.AllowedOperators)
			return utils.AppendPrefixOnSlice(remainder, keyOperatorCompletion), cobra.ShellCompDirectiveDefault
		}

		// Multiple key matched
		return utils.AppendPrefixOnSlice(remainder, keyMatches), cobra.ShellCompDirectiveDefault
	}

	// Writing Operator
	writingOp := key != utils.EmptyString && op != utils.EmptyString && value == utils.EmptyString
	if writingOp {
		cobra.CompDebug("writingOp case", true)
		opMatches := utils.Filter(ansible.AllowedOperators, func(s string) bool {
			return strings.HasPrefix(s, op)
		})

		// No matches return empty slice
		if len(opMatches) == 0 {
			return EmptyCompletion, cobra.ShellCompDirectiveDefault
		}

		// Match Op complete value
		if len(opMatches) == 1 {
			op := opMatches[0]
			// Add op, key and remainder
			return utils.AppendPrefixOnSlice(remainder, utils.AppendPrefixOnSlice(key, utils.AppendPrefixOnSlice(op, k8s.InventoryValues(key)))), cobra.ShellCompDirectiveDefault
		}

		// Multiple Op possible
		return utils.AppendPrefixOnSlice(remainder, utils.AppendPrefixOnSlice(key, opMatches)), cobra.ShellCompDirectiveDefault
	}

	// Writing value
	writingValue := key != utils.EmptyString && op != utils.EmptyString && value != utils.EmptyString
	if writingValue {
		cobra.CompDebug("writingValue case", true)
		matches := utils.Filter(k8s.InventoryValues(key), func(s string) bool {
			return strings.HasPrefix(s, value)
		})

		// No matches return empty slice
		if len(matches) == 0 {
			return EmptyCompletion, cobra.ShellCompDirectiveDefault
		}

		return utils.AppendPrefixOnSlice(remainder, utils.AppendPrefixOnSlice(key, utils.AppendPrefixOnSlice(op, matches))), cobra.ShellCompDirectiveDefault

	}

	cobra.CompDebug("no case match", true)
	return EmptyCompletion, cobra.ShellCompDirectiveDefault
}

func playbookCompletion(toComplete string, path string) ([]string, cobra.ShellCompDirective) {
	logrus.SetLevel(logrus.PanicLevel)
	k8s := ansible.Projects.LoadFromPath(path)
	return k8s.PlaybookPaths(), cobra.ShellCompDirectiveDefault
}

func tagsCompletion(toComplete string, path string, playbookPath string) ([]string, cobra.ShellCompDirective) {
	logrus.SetLevel(logrus.PanicLevel)
	var _ = regexp.MustCompile("([\\w-.\\/]+)([,]|)")

	if len(playbookPath) == 0 {
		return nil, cobra.ShellCompDirectiveDefault
	}
	project := ansible.Projects.LoadFromPath(path)

	playbook, err := project.PlaybookPath(playbookPath)

	if err != nil {
		cobra.CompDebug(err.Error(), true)
		return nil, cobra.ShellCompDirectiveDefault
	}

	return playbook.AllTags().List(), cobra.ShellCompDirectiveDefault
}
