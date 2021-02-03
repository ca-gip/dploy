package ansible

import (
	"fmt"
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/karrick/godirwalk"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var AnsibleValidRoleSubFolders = utils.NewSetFromSlice("tasks", "meta")

type Role struct {
	AbsolutePath string
	Name         string    `yaml:"role"`
	Tasks        []Task    `yaml:"tasks"`
	Tags         utils.Set `yaml:"tags"`
}

func (role *Role) AllTags() *utils.Set {
	tags := utils.NewSet()
	tags.Concat(role.Tags.List())

	if role.Tasks != nil {
		for _, task := range role.Tasks {
			tags.Concat(task.Tags.List())
		}
	}
	return tags
}

// Gather inventory files LoadFromPath a Parent directory
// Using a recursive scan. All non inventory files are ignored ( not .ini file )
// All sub parent directory added like label in the inventory
func (role *Role) LoadFromPath(rootPath string) (err error) {

	if role.Tasks == nil {
		role.Tasks = []Task{}
	}

	absRoot, err := filepath.Abs(rootPath + "/roles/" + role.Name)

	if err != nil {
		log.Debug("The role ", role.Name, "can't be read. Error:", err.Error())
		return
	}

	err = godirwalk.Walk(absRoot, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {

			if de.IsDir() && !AnsibleValidRoleSubFolders.Contains(de.Name()) {
				return nil
			}

			if !strings.Contains(filepath.Base(osPathname), ".yml") {
				return nil
			}

			binData, err := ioutil.ReadFile(osPathname)
			// IMPORTANT: Yaml and Json parser need a root element,
			// They can't read a raw list.
			content := fmt.Sprintf("tasks:\n%s", string(binData))

			if err != nil {
				log.Debug("Cannot read role file: ", osPathname, ". Error:", err.Error())
			}

			var roleTasks Role
			err = yaml.Unmarshal([]byte(content), &roleTasks)

			if err != nil {
				log.Warn("Error during role parsing: ", utils.WrapRed(osPathname), ". More info in trace level.")
				log.Trace("Err:", err.Error())
				return nil
			}

			if roleTasks.Tasks == nil {
				return nil
			}

			var resTask []Task
			for _, task := range roleTasks.Tasks {
				resTask = append(resTask, task)
			}
			role.Tasks = resTask
			log.Debug("Available tags for role ", utils.WrapGrey(osPathname), " are: ", role.AllTags().List())
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			log.Error(err)
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})
	return
}
