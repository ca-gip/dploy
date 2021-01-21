package ansible

import (
	"bufio"
	"fmt"
	"github.com/karrick/godirwalk"
	"github.com/relex/aini"
	"os"
	"path/filepath"
	"strings"
)

type Inventory struct {
	AbsolutePath string
	RootPath     *string
	PathTags     []string
	Data         *aini.InventoryData
}

func (i *Inventory) make() {
	if i == nil {
		return
	}

	if strings.Contains(filepath.Base(i.AbsolutePath), ".ini") {
		if file, err := os.Open(i.AbsolutePath); err == nil {
			reader := bufio.NewReader(file)
			if data, err := aini.Parse(reader); err == nil {
				i.Data = data
			}
		}
	}
}

func (i *Inventory) RelativePath() string {
	return strings.TrimPrefix(i.AbsolutePath, *i.RootPath+"/")
}

// Gather inventory files UnmarshallPath a Parent directory
// Using a recursive scan. All non inventory files are ignored ( not .ini file )
// All sub parent directory added like label in the inventory
func readInventories(rootPath string) (result []*Inventory, err error) {
	absRoot, err := filepath.Abs(rootPath)

	if err != nil {
		return
	}

	err = godirwalk.Walk(absRoot, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if strings.Contains(osPathname, "vars") || strings.Contains(osPathname, "template") {
				return godirwalk.SkipThis
			}

			if !strings.Contains(filepath.Base(osPathname), ".ini") {
				return nil
			}
			pathMetas := strings.Split(strings.TrimSuffix(strings.TrimPrefix(osPathname, absRoot), fmt.Sprintf("/%s", de.Name())), "/")

			result = append(result, &Inventory{
				AbsolutePath: osPathname,
				RootPath:     &rootPath,
				PathTags:     pathMetas,
			})
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})
	return
}
