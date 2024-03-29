package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ca-gip/dploy/internal/ansible"
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var EmptyCompletion []string

func extractMultipleCompletion(toComplete string) (remainder string, current string) {
	toCompletes := strings.Split(toComplete, ",")

	if len(toCompletes) == 1 {
		return utils.EmptyString, toCompletes[0]
	}

	remainder = fmt.Sprintf("%s%s", strings.Join(toCompletes[:len(toCompletes)-1], ","), ",")
	current = toCompletes[len(toCompletes)-1]
	return
}

// TODO : Update regex to allow var with hyphen
func filterCompletion(toComplete string, path string) ([]string, cobra.ShellCompDirective) {
	logrus.SetLevel(logrus.PanicLevel)

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
	project := ansible.Projects.LoadFromPath(path)
	return project.PlaybookPaths(), cobra.ShellCompDirectiveDefault
}

func tagsCompletion(toComplete string, path string, playbookPath string) ([]string, cobra.ShellCompDirective) {
	logrus.SetLevel(logrus.PanicLevel)

	if len(playbookPath) == 0 {
		return EmptyCompletion, cobra.ShellCompDirectiveDefault
	}

	project := ansible.Projects.LoadFromPath(path)
	playbook, err := project.PlaybookPath(playbookPath)

	if err != nil {
		cobra.CompDebug(err.Error(), true)
		return EmptyCompletion, cobra.ShellCompDirectiveDefault
	}

	remainder, current := extractMultipleCompletion(toComplete)
	cobra.CompDebug(fmt.Sprintf("extract muitple: remainder:%s current:%s\n", remainder, current), true)

	matches := utils.Filter(playbook.AllTags().List(), func(s string) bool {
		return strings.HasPrefix(s, current)
	})

	return utils.AppendPrefixOnSlice(remainder, matches), cobra.ShellCompDirectiveDefault
}

func hostPatternCompletion(toComplete string, path string) ([]string, cobra.ShellCompDirective) {
	logrus.SetLevel(logrus.PanicLevel)
	project := ansible.Projects.LoadFromPath(path)
	return append(project.InventoryHost(), project.InventoryGroups()...), cobra.ShellCompDirectiveDefault
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [Y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" || response == "" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func varListToMap(varsList []string) (map[string]interface{}, error) {
	vars := map[string]interface{}{}

	for _, v := range varsList {
		tokens := strings.Split(v, "=")

		if len(tokens) != 2 {
			return nil, errors.New(fmt.Sprintf("Invalid extra variable format on '%s'", v))
		}
		vars[tokens[0]] = tokens[1]
	}

	return vars, nil
}
