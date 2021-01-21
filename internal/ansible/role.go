package ansible

import (
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/karrick/godirwalk"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Role struct {
	AbsolutePath string
	Name         string    `yaml:"role"`
	Tasks        []Task    `yaml:"tasks"`
	Tags         utils.Set `yaml:"tags"`
}

func (role *Role) AllTags() (tags *utils.Set) {
	tags = utils.NewSet()
	for _, task := range role.Tasks {
		tags.Concat(task.Tags.List())
	}
	log.Debug("tags:::", tags.List())
	tags.Concat(role.Tags.List())
	return
}

// Gather inventory files LoadFromPath a Parent directory
// Using a recursive scan. All non inventory files are ignored ( not .ini file )
// All sub parent directory added like label in the inventory
func (role *Role) LoadFromPath(rootPath string) (err error) {
	absRoot, err := filepath.Abs(rootPath + "/roles/" + role.Name)
	log.Debug("reading ", role.Name, "at: ", absRoot)

	if err != nil {
		log.Debug("The role ", role.Name, "can't be read. Error:", err.Error())
		return
	}

	err = godirwalk.Walk(absRoot, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {

			if !strings.Contains(filepath.Base(osPathname), ".yml") {
				return nil
			}

			binData, err := ioutil.ReadFile(osPathname)
			if err != nil {
				log.Debug("Cannot read file: ", osPathname, ". Error:", err.Error())
			}

			var tasks []Task
			err = yaml.Unmarshal(binData, &tasks)

			if err != nil {
				log.Debug("Error reading role", osPathname, "err:", err.Error())
			} else {
				log.Debug("task is", tasks)
			}

			for _, task := range tasks {
				role.Tasks = append(role.Tasks, task)
			}

			log.Debug(osPathname, "tags in role tags:", role.AllTags())

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
