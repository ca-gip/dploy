package services

import (
	"log"
)

type Playbook struct {
	Name string
}

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
func LoadFromPath(rootPath string) (project Project) {
	project = Project{
		Name: rootPath,
	}
	inventories, err := readInventories(rootPath)
	if err != nil {
		log.Fatalln("Cannot parse directory: ", err.Error())
	}
	for _, inventory := range inventories {
		inventory.make()
	}
	project.Inventories = inventories
	return
}
