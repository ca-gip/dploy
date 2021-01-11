package services

import (
	"log"
)

type Project struct {
	Path        *string
	Inventories Inventories
	Playbooks   Playbooks
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
