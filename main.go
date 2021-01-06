package main

import (
	"fmt"
	"github.com/ca-gip/dploy/cmd"
	"github.com/ca-gip/dploy/internal/services"
	"os"
)

func main() {

	home, _ := os.UserHomeDir()
	path := fmt.Sprintf("%s/%s", home, "Projects/ansible-kube/inventories")
	k8s := services.LoadFromPath(path, path+"/..")

	filter := map[string]string{"customer": "cagip"}
	filteredInventories := k8s.FilterByVarsOr(filter)
	fmt.Println("Filtering ", len(filteredInventories), "/", len(k8s.Inventories))
	for _, i := range filteredInventories {
		fmt.Println(i.FilePath)
	}

	fmt.Println("Playbooks")

	for _, i := range k8s.Playbooks {
		fmt.Println(i.Name, i.Plays)
		for _, t := range i.Plays {
			fmt.Printf("\ntags:%v", t.Tags)
		}
	}

	cmd.Execute()
}
