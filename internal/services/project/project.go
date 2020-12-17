package project

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
	Name     string
	Path     string
	IsGroup  bool
	Children []*Inventory
	Parent   *Inventory
	Data     *aini.InventoryData
}

type Playbook struct {
	Name string
}

type Project struct {
	Name        string
	Inventories []*Inventory
	Playbooks   []*Playbook
}

func NewProject(root string) (result *Inventory, err error) {
	absRoot, err := filepath.Abs(root)

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

func MarkInventoryGroup(root *Inventory) {
	if root == nil {
		return
	}

	if len(root.Children) <= 1 {
		root.IsGroup = false
	} else {
		root.IsGroup = true
	}

	for _, child := range root.Children {
		MarkInventoryGroup(child)
	}

	return

}

func ParseInventory(root *Inventory) {
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
		ParseInventory(child)
	}
}

func (i *Inventory) FilterByVarsOr(filters map[string]string) (filtered *[]Inventory) {
	acc := &[]Inventory{}
	return FilterByVarsOr(i, filters, acc)
}

func FilterByVarsOr(root *Inventory, filters map[string]string, acc *[]Inventory) (filtered *[]Inventory) {
	if root == nil {
		return
	}

	if root.Data != nil {
		for key, value := range filters {
			if root.Data.Groups["all"].Vars[key] == value {
				*acc = append(*acc, *root)
			}
		}
	}

	for _, child := range root.Children {
		FilterByVarsOr(child, filters, acc)
	}

	return acc
}

func (i *Inventory) ReadInventory() {
	readInventory(i)
}

func readInventory(root *Inventory) {
	if root == nil {
		return
	}

	fmt.Printf("%+v\n", *root)

	for _, child := range root.Children {
		readInventory(child)
	}
}
