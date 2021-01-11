package services

import (
	"log"
)

type Project struct {
	Path        *string
	Inventories []*Inventory
	Playbooks   []*Playbook
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

func (project *Project) FilterFromVars(filters []Filter) (filtered []*Inventory) {

	if len(filters) == 0 {
		return project.Inventories
	}

	for _, inventory := range project.Inventories {
		if inventory.Data != nil {

			type condition = string
			matchFilter := make(map[condition]bool)

			for _, filter := range filters {
				if filter.Eval(inventory.Data.Groups["all"].Vars[filter.Key]) {
					matchFilter[filter.GetRaw()] = true
				} else {
					matchFilter[filter.GetRaw()] = false
				}
			}

			if AllTrue(matchFilter) {
				filtered = append(filtered, inventory)
			}

		}
	}
	return
}

// TODO : https://www.davidkaya.com/sets-in-golang/ ?
func (project *Project) GetInventoryKeys() (keys []string) {
	var Set = struct{}{}
	uniqueKeys := make(map[string]struct{})
	for _, inventory := range project.Inventories {
		if inventory.Data != nil {
			for key, _ := range inventory.Data.Groups["all"].Vars {
				uniqueKeys[key] = Set
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

func (project *Project) GetPlaybooks() (values []string) {
	for _, playbook := range project.Playbooks {
		values = append(values, playbook.RelativePath())
	}

	return
}

func (project *Project) GetPlaybook(path string) *Playbook {
	for _, playbook := range project.Playbooks {
		if path == playbook.RelativePath() {
			return playbook
		}
	}
	return nil
}

func (project *Project) GetInventoryLength() int {
	return len(project.Inventories)
}

// TODO: Add assert on file system ( readable, permissions ...)
func LoadFromPath(projectDirectory string) (project Project) {
	project = Project{
		Path:      &projectDirectory,
		Playbooks: nil,
	}
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
