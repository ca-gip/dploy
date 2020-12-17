package main

import (
	"fmt"
	"github.com/ca-gip/dploy/internal/services/project"
	"os"
)

func main() {

	home, _ := os.UserHomeDir()
	path := fmt.Sprintf("%s/%s", home, "Projects/ansible-kube/inventories")

	k8s := project.NewInventory(path)

	filter := map[string]string{"customer": "cagip", "network_name_suffix": "k8s"}
	fmt.Print(k8s.FilterByVarsOr(filter))

}
