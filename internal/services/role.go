package services

import (
	"fmt"
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/ghodss/yaml"
	"github.com/karrick/godirwalk"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Role struct {
	AbsolutePath string
	Name         string `yaml:"role"`
	Tasks        []Tasks
	Tags         []string `yaml:"tags,omitempty" yaml:"tags,omitempty"`
}

// Gather inventory files from a Parent directory
// Using a recursive scan. All non inventory files are ignored ( not .ini file )
// All sub parent directory added like label in the inventory
func (role *Role) ReadRole(rootPath string, pathTags ...string) (err error) {
	absRoot, err := filepath.Abs(rootPath + "/roles/" + role.Name)

	if err != nil {
		// TODO log erreor reading role
		return
	}

	err = godirwalk.Walk(absRoot, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {

			tags := utils.NewSet()

			if !strings.Contains(filepath.Base(osPathname), ".yml") {
				return nil
			}

			binData, err := ioutil.ReadFile(osPathname)
			if err != nil {
				// TODO log fatal, unable to read file
				fmt.Println(err)
			}

			var tasks []Task
			err = yaml.Unmarshal([]byte(binData), &tasks)
			for _, task := range tasks {
				tags.Concat(task.Tags)
			}

			tasks = append(tasks, Task{Tags: tags.List()})
			if len(tags.List()) > 0 {
				fmt.Println("Task tags:", tags.List())
			}
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})
	return
}
