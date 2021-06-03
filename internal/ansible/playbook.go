package ansible

import (
	"errors"
	"fmt"
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/karrick/godirwalk"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type playbooks struct{}

var Playbooks = playbooks{}

type Playbook struct {
	absolutePath string
	rootPath     *string
	Plays        []*Play
}

func (playbook *Playbook) AllTags() (tags *utils.Set) {
	tags = utils.NewSet()
	for _, play := range playbook.Plays {
		tags.Concat(play.AllTags().List())
		for _, role := range play.Roles {
			tags.Concat(role.AllTags().List())
		}
	}

	return
}

func (playbook *Playbook) RelativePath() string {
	return strings.TrimPrefix(playbook.absolutePath, *playbook.rootPath+"/")
}

// Unmarshall a playbook from file
func (p *playbooks) unmarshallFromPath(playbookPath string, rootPath string) (playbook *Playbook, err error) {
	// Try to check playbook content
	binData, err := ioutil.ReadFile(playbookPath)

	// IMPORTANT: Yaml and Json parser need a root element,
	// They can't read a raw list.
	content := fmt.Sprintf("plays:\n%s", string(binData))

	if err != nil {
		log.Error("Cannot read playbook", playbookPath, ". Error: ", err.Error())
		return nil, err
	}
	err = yaml.Unmarshal([]byte(content), &playbook)
	if err != nil {
		log.Debug("Skip ", playbookPath, " not a playbook ", err.Error())
		return nil, err
	}
	if len(playbook.Plays) == 0 || playbook.Plays[0].Hosts == utils.EmptyString {
		return nil, errors.New(fmt.Sprint("No play found inside the playbook: ", playbookPath))
	}

	for _, play := range playbook.Plays {
		for _, role := range play.Roles {
			role.LoadFromPath(rootPath)
		}
	}
	return
}

// Gather playbook files UnmarshallPath a Parent directory
// Using a recursive scan. All non playbook files are ignored ( not .yaml or .yml file )
func (p *playbooks) LoadFromPath(rootPath string) (result []*Playbook, err error) {
	absRoot, err := filepath.Abs(rootPath)

	if err != nil {
		log.Error(err)
		return
	}

	// Merge Play, Role and Task Tags for a playbook
	//allTags := utils.NewSet()

	err = godirwalk.Walk(absRoot, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if strings.Contains(de.Name(), "vars") || de.Name() == "template" || de.Name() == "roles" {
				return godirwalk.SkipThis
			}

			if !strings.Contains(filepath.Base(osPathname), ".yaml") && !strings.Contains(filepath.Base(osPathname), ".yml") {
				return nil
			}

			// Try to check playbook content
			playbook, err := Playbooks.unmarshallFromPath(osPathname, rootPath)

			if err != nil {
				return nil
			}

			// Browse Role Tags

			playbook.absolutePath = osPathname
			playbook.rootPath = &rootPath

			result = append(result, playbook)
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})

	return
}
