package main

import (
	"fmt"
	"github.com/ca-gip/dploy/internal/services/project"
	"os"
)

func main() {

	home, _ := os.UserHomeDir()
	path := fmt.Sprintf("%s/%s", home, "Projects/ansible-kube/inventories")

	fmt.Print(path)

	k8s, _ := project.NewProject(path)

	project.MarkInventoryGroup(k8s)
	project.ParseInventory(k8s)

	filter := map[string]string{"customer": "cagip", "network_name_suffix": "k8s"}
	fmt.Print(k8s.FilterByVarsOr(filter))

}
