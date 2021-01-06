package main

import (
	"fmt"
	"github.com/ca-gip/dploy/internal/services"
	"os"
)

func main() {

	home, _ := os.UserHomeDir()
	path := fmt.Sprintf("%s/%s", home, "Projects/ansible-kube/inventories/")

	k8s := services.LoadFromPath(path)

	filter := map[string]string{"customer": "cagip"}
	filteredInventories := k8s.FilterByVarsOr(filter)
	fmt.Println("Filtering ", len(filteredInventories), "/", len(k8s.Inventories))
	for _, i := range filteredInventories {
		fmt.Println(i.FilePath)
	}

}
