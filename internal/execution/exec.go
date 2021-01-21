package execution

import (
	ansibler "github.com/apenella/go-ansible"
	"github.com/ca-gip/dploy/internal/ansible"
	"strings"
)

func Exec(command ansible.Command) (err error) {
	for _, playbook := range MakeCommands(command) {
		ansibler.AnsibleForceColor()
		err = playbook.Run()
		if err != nil {
			return
		}
	}
	return
}

func MakeCommands(command ansible.Command) (commands []ansibler.AnsiblePlaybookCmd) {
	for _, inventory := range command.Inventory {
		ansiblePlaybookOptions := &ansibler.AnsiblePlaybookOptions{
			Inventory:         inventory.RelativePath(),
			Limit:             strings.Join(command.Limit, ","),
			Tags:              strings.Join(command.Tags, ","),
			VaultPasswordFile: command.VaultPasswordFile,
		}
		playbook := ansibler.AnsiblePlaybookCmd{
			Playbook: command.Playbook.RelativePath(),
			Options:  ansiblePlaybookOptions,
		}
		commands = append(commands, playbook)
	}
	return
}
