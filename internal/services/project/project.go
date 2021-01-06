package project

type Playbook struct {
	Name string
}

type Project struct {
	Name        string
	Inventories []*InventoryPath
	Playbooks   []*Playbook
}

func (project *Project) FilterByVarsOr(filters map[string]string) (filtered []*InventoryPath) {
	if project == nil {
		return
	}

	for _, inventory := range project.Inventories {
		for key, value := range filters {
			if inventory.Data.Groups["all"].Vars[key] == value {
				filtered = append(filtered, inventory)
			}
		}
		if inventory.Data != nil {
			for key, value := range filters {
				if inventory.Data.Groups["all"].Vars[key] == value {
					filtered = append(filtered, inventory)
				}
			}
		}
	}
	return
}
