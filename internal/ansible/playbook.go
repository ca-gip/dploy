package ansible

import (
	"fmt"
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/karrick/godirwalk"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Playbook struct {
	absolutePath string
	rootPath     *string
	Plays        []Play
}

const decoderTagName = "tags"

func (playbook *Playbook) AllTags() (tags *utils.Set) {
	tags = utils.NewSet()
	for _, play := range playbook.Plays {
		tags.Concat(play.AllTags().List())
	}
	return
}

func (playbook *Playbook) RelativePath() string {
	return strings.TrimPrefix(playbook.absolutePath, *playbook.rootPath+"/")
}

func ReadFromFile(osPathname string) (playbook Playbook) {
	// Try to check playbook content
	binData, err := ioutil.ReadFile(osPathname)

	// IMPORTANT: Yaml and Json parser need a root element,
	// They can't read a raw list.
	content := fmt.Sprintf("plays:\n%s", string(binData))

	if err != nil {
		fmt.Println("Cannot read playbook", osPathname, ". Error: ", err.Error())
		return
	}
	err = yaml.Unmarshal([]byte(content), &playbook)
	if err != nil {
		fmt.Println("Skip", osPathname, " not a playbook ", err.Error())
		return
	}
	if len(playbook.Plays) == 0 {
		fmt.Println("No play found inside the playbook: ", osPathname)
		return
	}
	if playbook.Plays[0].Hosts == utils.EmptyString {
		fmt.Println("No play found inside the playbook: ", osPathname)
		return
	}
	return
}

// Gather playbook files from a Parent directory
// Using a recursive scan. All non playbook files are ignored ( not .yaml or .yml file )
func readPlaybook(rootPath string) (result []*Playbook, err error) {
	absRoot, err := filepath.Abs(rootPath)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("reading playbook")
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
			playbook := ReadFromFile(osPathname)

			// Browse Role Tags
			for _, play := range playbook.Plays {

				allTags.Concat(play.AllTags().List())
				fmt.Println("Play tags are: ", play.Tags)
				for _, role := range play.Roles {
					role.ReadRoleTasks(rootPath)
					fmt.Println("  Role info", role.AllTags())
					allTags.Concat(role.AllTags().List())
				}
			}

			playbook.absolutePath = osPathname
			playbook.rootPath = &rootPath

			result = append(result, &playbook)
			fmt.Println("Available tags are :", playbook.AllTags())
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})

	return
}
