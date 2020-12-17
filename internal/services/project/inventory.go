package project

import "fmt"

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
