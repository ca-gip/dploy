package services

import (
	"fmt"
	"log"
)

type Project struct {
	Name        string
	Inventories []*Inventory
	Playbooks   []*Playbook
}

func (project *Project) FilterByVarsOr(filters map[string]string) (filtered []*Inventory) {
	if project == nil {
		return
	}

	for _, inventory := range project.Inventories {
		if inventory.Data != nil {
			for key, value := range filters {
				if inventory.Data.Groups["all"].Vars[key] == value {
					filtered = append(filtered, inventory)
				}
			}
		}
	}
	return
}

// TODO: Add assert on file system ( readable, permissions ...)
func LoadFromPath(inventoryPath string, playbookPath string) (project Project) {
	project = Project{
		Name:      inventoryPath,
		Playbooks: nil,
	}
	fmt.Println(playbookPath)
	playbooks, errPlaybooks := readPlaybook(playbookPath)
	inventories, errInventories := readInventories(inventoryPath)
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
