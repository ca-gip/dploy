package ansible

import (
	"errors"
	"fmt"
	"github.com/ca-gip/dploy/internal/utils"
	log "github.com/sirupsen/logrus"
)

type projects struct{}

var Projects = projects{}

type Project struct {
	Path        *string
	Inventories []*Inventory
	Playbooks   []*Playbook
}

// Returns inventory that match all the conditions
func (p Project) FilterInventory(filters []Filter) (filtered []*Inventory) {
	if len(filters) == 0 {
		return
	}

	for _, inventory := range p.Inventories {
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

			if utils.MapHasAllTrue(matchFilter) {
				filtered = append(filtered, inventory)
			}

		}
	}
	return
}

// Returns variables name contained in the all section of all inventory files
func (p Project) InventoryKeys() (keys []string) {
	keySet := utils.NewSet()
	for _, inventory := range p.Inventories {
		if inventory.Data != nil {
			for key := range inventory.Data.Groups["all"].Vars {
				keySet.Add(key)
			}
		}
	}
	return keySet.List()
}

// Given a variable name returns value across multiples inventory files
func (p Project) InventoryValues(key string) (values []string) {
	valueSet := utils.NewSet()
	for _, inventory := range p.Inventories {
		if inventory.Data != nil {
			if value := inventory.Data.Groups["all"].Vars[key]; value != "" {
				valueSet.Add(value)
			}
		}
	}
	return valueSet.List()
}

func (p Project) PlaybookPaths() (values []string) {
	for _, playbook := range p.Playbooks {
		values = append(values, playbook.RelativePath())
	}
	return
}

func (p Project) PlaybookPath(path string) (playbook *Playbook, err error) {
	for _, playbook := range p.Playbooks {
		if path == playbook.RelativePath() {
			return playbook, nil
		}
	}
	return nil, errors.New(fmt.Sprint("No playbook found at path: ", path))
}

// TODO: Add assert on file system ( readable, permissions ...)
func (projects *projects) LoadFromPath(projectDirectory string) (project Project) {

	project = Project{
		Path:      &projectDirectory,
		Playbooks: nil,
	}

	playbooks, errPlaybooks := Playbooks.LoadFromPath(projectDirectory)
	inventories, errInventories := Inventories.LoadFromPath(projectDirectory)
	project.Playbooks = playbooks
	project.Inventories = inventories

	if errPlaybooks != nil {
		log.Fatalln("Cannot parse directory for playbooks: ", errPlaybooks.Error())
	}
	if errInventories != nil {
		log.Fatalln("Cannot parse directory for inventories: ", errInventories.Error())
	}
	return
}
