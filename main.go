package main

import (
	"github.com/ca-gip/dploy/cmd"
)

func main() {

	//home, _ := os.UserHomeDir()
	//path := fmt.Sprintf("%s/%s", home, "Projects/ansible-kube")
	//k8s := services.LoadFromPath(path)
	//fmt.Println("Filtering ", len(k8s.Inventories), "/", len(k8s.Inventories))
	//
	//fmt.Println("Playbooks")
	//
	//tpl := services.AnsibleCommandTpl{
	//	Inventory:         k8s.Inventories,
	//	Playbook:          k8s.Playbooks[1],
	//	Tags:              []string{"tag1", "tag2"},
	//	Limit:             []string{"limit1,limit2"},
	//	SkipTags:          []string{"testt"},
	//	Check:             true,
	//	Diff:              true,
	//	VaultPasswordFile: "/path/to/passwordfile",
	//	AskVaultPass:  false,
	//}
	//tpl.GenerateCmd()

	cmd.Execute()
}
