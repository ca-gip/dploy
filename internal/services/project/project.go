package project

type Playbook struct {
	Name string
}

type Project struct {
	Name        string
	Inventories []*Inventory
	Playbooks   []*Playbook
}
