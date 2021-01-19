package ansible

import (
	"fmt"
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/ghodss/yaml"
	"github.com/karrick/godirwalk"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Tasks struct {
	Tags []string `yaml:"tags,omitempty" yaml:"tags,omitempty"`
}

type Playbook struct {
	AbsolutePath string
	RootPath     *string
	Plays        []Play
	AllTags      utils.Set
}

func (playbook *Playbook) RelativePath() string {
	return strings.TrimPrefix(playbook.AbsolutePath, *playbook.RootPath+"/")
}

// Gather playbook files from a Parent directory
// Using a recursive scan. All non playbook files are ignored ( not .yaml or .yml file )
func readPlaybook(rootPath string) (result []*Playbook, err error) {
	absRoot, err := filepath.Abs(rootPath)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Merge Play, Role and Task Tags for a playbook
	allTags := utils.NewSet()

	err = godirwalk.Walk(absRoot, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if strings.Contains(osPathname, "vars") || strings.Contains(osPathname, "template") {
				return godirwalk.SkipThis
			}

			if !strings.Contains(filepath.Base(osPathname), ".yaml") && !strings.Contains(filepath.Base(osPathname), ".yml") {
				return nil
			}

			// Try to check playbook content
			var plays []Play
			binData, err := ioutil.ReadFile(osPathname)
			if err != nil {
				log.Error("Cannot read playbook", osPathname, ". Error: ", err.Error())
				return nil
			}
			err = yaml.Unmarshal([]byte(binData), &plays)
			if err != nil {
				log.Error("Cannot unmashal playbook data", osPathname, ". Error: ", err.Error())
				return nil
			}
			if plays == nil || len(plays) == 0 {
				log.Debug("No play found inside the playbook: ", osPathname)
				return nil
			}
			if plays[0].Hosts == utils.EmptyString {
				log.Debug("No play found inside the playbook: ", osPathname)
				return nil
			}

			// Browse Role Tags
			for _, play := range plays {
				allTags.Concat(play.Tags)
				for _, role := range play.Roles {
					role.ReadRole(rootPath)
					allTags.Concat(role.Tags)
				}
			}

			playbook := Playbook{
				RootPath:     &rootPath,
				AbsolutePath: osPathname,
				Plays:        plays,
				AllTags:      *allTags,
			}
			result = append(result, &playbook)
			log.Debug("Available tags are :", playbook.AllTags)
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})

	return
}
