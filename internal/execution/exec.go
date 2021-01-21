package execution

import (
	"bytes"
	"fmt"
	ansibler "github.com/apenella/go-ansible"
	"github.com/apenella/go-ansible/stdoutcallback/results"
	"github.com/ca-gip/dploy/internal/ansible"
	"strings"
)

func Exec(command ansible.Command) {
	var err error
	res := &results.AnsiblePlaybookJSONResults{}
	buff := new(bytes.Buffer)

	for _, playbook := range MakeCommands(command) {
		err = playbook.Run()
		if err != nil {
			fmt.Println(err.Error())
		}

		res, err = results.JSONParse(buff.Bytes())

		if err != nil {
			fmt.Printf("there is an error %s", err)
		}

		fmt.Println(res.String())
	}
}

func MakeCommands(command ansible.Command) (commands []ansibler.AnsiblePlaybookCmd) {
	for _, inventory := range command.Inventory {
		ansiblePlaybookOptions := &ansibler.AnsiblePlaybookOptions{
			Inventory:         inventory.AbsolutePath,
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
