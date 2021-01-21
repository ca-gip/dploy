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
	tags.Concat(role.Tags.List())
	return
}

// Gather inventory files from a Parent directory
// Using a recursive scan. All non inventory files are ignored ( not .ini file )
// All sub parent directory added like label in the inventory
func (role *Role) ReadRoleTasks(rootPath string, pathTags ...string) (err error) {
	absRoot, err := filepath.Abs(rootPath + "/roles/" + role.Name)
	fmt.Println("reading ", role.Name, "at: ", absRoot)

	if err != nil {
		fmt.Println("The role ", role.Name, "can't be read. Error:", err.Error())
		return
	}

	fmt.Println(role.Name)
	err = godirwalk.Walk(absRoot, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {

			if !strings.Contains(filepath.Base(osPathname), ".yml") {
				return nil
			}

			binData, err := ioutil.ReadFile(osPathname)
			if err != nil {
				fmt.Println("Cannot read file: ", osPathname, ". Error:", err.Error())
			}

			var tasks []Task
			err = yaml.Unmarshal([]byte(binData), &tasks)

			if err != nil {
				fmt.Println("Error reading role", osPathname, "err:", err.Error())
			} else {
				fmt.Println("Task is", tasks)
			}

			for _, task := range tasks {
				role.Tasks = append(role.Tasks, task)
			}
			fmt.Println(osPathname, "tags in role tags:", role.AllTags())

			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			fmt.Println(err.Error())
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})
	return
}
