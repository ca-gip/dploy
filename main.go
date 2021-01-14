package main

import (
	"github.com/ca-gip/dploy/cmd"
	"github.com/ca-gip/dploy/internal/execution"
)

func main() {

	//home, _ := os.UserHomeDir()
	//path := fmt.Sprintf("%s/%s", home, "Projects/ansible-kube")
	//k8s := ansible.LoadFromPath(path)
	//
	//filter := []string{"customer!=cacf_corp_hors_prod"}
	//filteredInventories := k8s.FilterFromVars(filter)
	//fmt.Println("Filtering ", len(filteredInventories), "/", len(k8s.Inventories))
	//for _, i := range filteredInventories {
	//	fmt.Println(i.AbsolutePath)
	//}
	//
	//fmt.Println("Playbooks")
	//
	//tpl := ansible.AnsibleCommandTpl{
	//	Inventory:         filteredInventories,
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
	execution.Test()

	cmd.Execute()
}
