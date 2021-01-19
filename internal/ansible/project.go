package ansible

import (
	log "github.com/sirupsen/logrus"
)

type Project struct {
	Path        *string
	Inventories Inventories
	Playbooks   Playbooks
}

// TODO: Add assert on file system ( readable, permissions ...)
func LoadFromPath(projectDirectory string) (project Project) {

	log.Info("info")
	log.Debug("debug")
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
