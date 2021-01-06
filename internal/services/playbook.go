package services

import (
	"fmt"
	"github.com/karrick/godirwalk"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Role struct {
	Role string      `yaml:"role"`
	Tags interface{} `yaml:"tags,omitempty" yaml:"tags,omitempty"`
}

type Play struct {
	Hosts string      `yaml:"hosts"`
	Roles []Role      `yaml:"roles"`
	Tags  interface{} `yaml:"tags,omitempty"`
}

type Playbook struct {
	Name  string
	Plays []Play
}

// Gather playbook files from a Parent directory
// Using a recursive scan. All non playbook files are ignored ( not .yaml or .yml file )
func readPlaybook(rootPath string) (result []*Playbook, err error) {
	absRoot, err := filepath.Abs(rootPath)

	if err != nil {
		fmt.Println(err)
		return
	}

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
				// TODO add debug for read
				return nil
			}
			err = yaml.Unmarshal([]byte(binData), &plays)
			if err != nil {
				// TODO add debug for unmarshaling
				return nil
			}

			if plays == nil || len(plays) == 0 {
				// TODO Log debug no play found
				return nil
			}

			if plays[0].Hosts == "" {
				// TODO Log debug do not seems to be a playbook
				return nil
			}

			result = append(result, &Playbook{Name: osPathname, Plays: plays})
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})
	return
}
