package ansible

type Playbooks []*Playbook

func (playbooks Playbooks) GetPlaybooks() (values []string) {
	for _, playbook := range playbooks {
		values = append(values, playbook.RelativePath())
	}
	return
}

func (playbooks Playbooks) GetPlaybook(path string) *Playbook {
	for _, playbook := range playbooks {
		if path == playbook.RelativePath() {
			return playbook
		}
	}
	return nil
}
