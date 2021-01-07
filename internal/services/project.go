package services

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Project struct {
	Path        *string
	Inventories []*Inventory
	Playbooks   []*Playbook
}

var filterRegex = regexp.MustCompile("(\\w*)(\\W*)(\\w*)")

func ParseFilter(filter string) (key, op, value string) {
	result := filterRegex.FindStringSubmatch(filter)
	return result[1], result[2], result[3]
}

func ConditionEval(left, right, op string) bool {
	switch op {
	case "==":
		return strings.EqualFold(left, right)
	case "!=":
		return !strings.EqualFold(left, right)
	case "$=":
		return strings.HasSuffix(left, right)
	case "~=":
		return strings.Contains(left, right)
	case "^=":
		return strings.HasPrefix(left, right)
	default:
		log.Fatalf("Unsuported filter operation %s", op)
		return false
	}
}

func AllTrue(a map[string]bool) bool {
	if len(a) == 0 {
		return false
	}

	for _, value := range a {
		if !value {
			return false
		}
	}

	return true
}

func (project *Project) FilterFromVars(filters []string) (filtered []*Inventory) {
	for _, inventory := range project.Inventories {
		if inventory.Data != nil {

			type condition = string
			matchFilter := make(map[condition]bool)

			for _, filter := range filters {
				key, op, value := ParseFilter(filter)
				inventoryValue := inventory.Data.Groups["all"].Vars[key]
				if ConditionEval(inventoryValue, value, op) {
					matchFilter[filter] = true
				} else {
					matchFilter[filter] = false
				}
			}

			if AllTrue(matchFilter) {
				filtered = append(filtered, inventory)
			}

		}
	}
	return
}

func (project *Project) GetInventoryKeys() (keys []string) {
	uniqueKeys := make(map[string]interface{})
	for _, inventory := range project.Inventories {
		if inventory.Data != nil {
			for key, _ := range inventory.Data.Groups["all"].Vars {
				uniqueKeys[key] = nil
			}
		}
	}

	for key, _ := range uniqueKeys {
		keys = append(keys, key)
	}

	return
}

func (project *Project) GetInventoryValues(key string) (values []string) {
	uniqueValues := make(map[string]interface{})
	for _, inventory := range project.Inventories {
		if inventory.Data != nil {
			value := inventory.Data.Groups["all"].Vars[key]
			uniqueValues[value] = nil
		}
	}

	for key, _ := range uniqueValues {
		values = append(values, key)
	}

	return
}

// TODO: Add assert on file system ( readable, permissions ...)
func LoadFromPath(projectDirectory string) (project Project) {
	project = Project{
		Path:      &projectDirectory,
		Playbooks: nil,
	}
	fmt.Println(projectDirectory)
	playbooks, errPlaybooks := readPlaybook(projectDirectory)
	inventories, errInventories := readInventories(projectDirectory)
	project.Playbooks = playbooks
	project.Inventories = inventories

	if errPlaybooks != nil {
		log.Fatalln("Cannot parse directory for playbooks: ", errPlaybooks.Error())
	}
	if errInventories != nil {
		log.Fatalln("Cannot parse directory for inventories: ", errInventories.Error())
	}
	for _, inventory := range inventories {
		inventory.make()
	}
	return
}
