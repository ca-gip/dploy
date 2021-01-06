package project

import "fmt"

func (i *InventoryPath) ReadInventory() {
	readInventory(i)
}

func readInventory(root *InventoryPath) {
	if root == nil {
		return
	}

	fmt.Printf("%+v\n", *root)
}
