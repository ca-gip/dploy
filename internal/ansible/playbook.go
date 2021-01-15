package ansible

import (
	"fmt"
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/karrick/godirwalk"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog/v2"
	"path/filepath"
	"strings"
)

type Tasks struct {
	Tags interface{} `yaml:"tags,omitempty"`
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
			if strings.Contains(de.Name(), "vars") || de.Name() == "template" || de.Name() == "roles" {
				return godirwalk.SkipThis
			}

			if !strings.Contains(filepath.Base(osPathname), ".yaml") && !strings.Contains(filepath.Base(osPathname), ".yml") {
				return nil
			}

			// Try to check playbook content
			var plays []Play
			binData, err := ioutil.ReadFile(osPathname)
			if err != nil {
				klog.Error("Cannot read playbook", osPathname, ". Error: ", err.Error())
				return nil
			}
			err = yaml.Unmarshal([]byte(binData), &plays)
			if err != nil {
				klog.Error("Skip", osPathname, " not an inventory ")
				return nil
			}
			if plays == nil || len(plays) == 0 {
				klog.Info("No play found inside the playbook: ", osPathname)
				return nil
			}
			if plays[0].Hosts == utils.EmptyString {
				klog.V(8).Info("No play found inside the playbook: ", osPathname)
				return nil
			}

			// Browse Role Tags
			for _, play := range plays {

				allTags.Concat(play.Tags())
				fmt.Println("Play tags are: ", play.Tags())
				for _, role := range play.Roles {
					role.ReadRole(rootPath)
					allTags.Concat(role.Tags())
				}
			}

			playbook := Playbook{
				RootPath:     &rootPath,
				AbsolutePath: osPathname,
				Plays:        plays,
				AllTags:      *allTags,
			}
			result = append(result, &playbook)
			klog.V(8).Info("Available tags are :", playbook.AllTags)
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})

	return
}
