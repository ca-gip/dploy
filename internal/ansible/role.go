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

type Role struct {
	AbsolutePath string
	Name         string `yaml:"role"`
	Tasks        []Tasks
	rawTags      interface{} `yaml:"tags,omitempty"`
}

func (role *Role) Tags() []string {
	return utils.InferSlice(role.rawTags)
}

// Gather inventory files from a Parent directory
// Using a recursive scan. All non inventory files are ignored ( not .ini file )
// All sub parent directory added like label in the inventory
func (role *Role) ReadRole(rootPath string, pathTags ...string) (err error) {
	absRoot, err := filepath.Abs(rootPath + "/roles/" + role.Name)

	if err != nil {
		klog.Error("The role ", role.Name, "can't be read. Error:", err.Error())
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
				klog.Error("Cannot read file: ", osPathname, ". Error:", err.Error())
			}

			var tasks []Task
			err = yaml.Unmarshal([]byte(binData), &tasks)
			for _, task := range tasks {
				fmt.Println("task t:", task.Tags())
				tags.Concat(task.Tags())
			}

			tasks = append(tasks, Task{rawTags: tags.List()})
			if len(tags.List()) > 0 {
				klog.Info("Task tags:", tags.List())
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
