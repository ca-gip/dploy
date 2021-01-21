package ansible

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type projects struct{}

var Projects = projects{}

type Project struct {
	Path        *string
	Inventories Inventories
	Playbooks   []Playbook
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
			return &playbook, nil
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
