package main

import (
	"fmt"
	"github.com/ca-gip/dploy/internal/services/project"
	"os"
)

func main() {

	home, _ := os.UserHomeDir()
	path := fmt.Sprintf("%s/%s", home, "Projects/ansible-kube/inventories/")

	k8s := project.NewInventory(path)

	filter := map[string]string{"customer": "cagip"}
	fmt.Println("Filtering ", len(filter), "/", len(k8s.Inventories))
	for _, i := range k8s.FilterByVarsOr(filter) {
		fmt.Println(i.FilePath, "of", len(k8s.Inventories))
	}

}
