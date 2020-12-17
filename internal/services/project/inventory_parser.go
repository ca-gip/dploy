package project

import (
	"bufio"
	"github.com/karrick/godirwalk"
	"github.com/relex/aini"
	"os"
	"path/filepath"
	"strings"
)

type Inventory struct {
	Name     string
	Path     string
	IsGroup  bool
	Children []*Inventory
	Parent   *Inventory
	Data     *aini.InventoryData
}

func scanInventoryDir(rootPath string) (result *Inventory, err error) {
	absRoot, err := filepath.Abs(rootPath)

	if err != nil {
		return
	}

	parents := make(map[string]*Inventory)

	err = godirwalk.Walk(absRoot, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if strings.Contains(osPathname, "vars") || strings.Contains(osPathname, "template") {
				return godirwalk.SkipThis
			}

			parents[osPathname] = &Inventory{
				Name:     filepath.Base(osPathname),
				Path:     osPathname,
				Children: make([]*Inventory, 0),
			}

			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})

	for path, inventory := range parents {
		parent, _ := parents[filepath.Dir(path)]
		if strings.EqualFold(path, absRoot) {
			result = inventory
		} else {
			inventory.Parent = parent
			parent.Children = append(parent.Children, inventory)
		}

	}

	return
}

func markInventoryGroup(root *Inventory) {
	if root == nil {
		return
	}

	if len(root.Children) <= 1 {
		root.IsGroup = false
	} else {
		root.IsGroup = true
	}

	for _, child := range root.Children {
		markInventoryGroup(child)
	}

	return

}

func parseInventory(root *Inventory) {
	if root == nil {
		return
	}

	if strings.Contains(filepath.Base(root.Path), ".ini") {
		if file, err := os.Open(root.Path); err == nil {
			reader := bufio.NewReader(file)
			if data, err := aini.Parse(reader); err == nil {
				root.Data = data
			}
		}
	}

	for _, child := range root.Children {
		parseInventory(child)
	}
}

func NewInventory(rootPath string) (inventory *Inventory) {
	inventory, _ = scanInventoryDir(rootPath)
	markInventoryGroup(inventory)
	parseInventory(inventory)
	return
}
