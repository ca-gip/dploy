package services

import (
	"os"
	"text/template"
)

type AnsibleCommandTpl struct {
	Inventory         []*Inventory
	Playbook          *Playbook
	Tags              []string
	Limit             []string
	SkipTags          []string
	Check             bool
	Diff              bool
	VaultPasswordFile string
	AskVaultPass      bool
}

var templateBash = `
# Dploy command result

{{- range $inventory := .Inventory }}
ansible-playbook -i {{ $inventory.RelativePath }} {{ $.Playbook.RelativePath }} 
{{- if $.Tags}} -t {{ range $i,$tag := $.Tags }}{{if gt $i 0 }},{{end}}{{ $tag }}{{ end }}{{- end}}
{{- if $.Limit}} -l {{ range $i,$limit := $.Limit }}{{if gt $i 0 }},{{end}}{{ $limit }}{{ end }}{{- end}}
{{- if $.Check}} -c{{- end}}
{{- if $.Diff}} --diff{{- end}}
{{- if $.AskVaultPass}} --ask-vault-password{{- end}}
{{- if $.VaultPasswordFile}} --vault-password-file {{ $.VaultPasswordFile}}{{- end}}
{{- end }}
`

func (tpl *AnsibleCommandTpl) GenerateCmd() {
	tmpl, _ := template.New("test").Parse(templateBash)
	tmpl.Execute(os.Stdout, tpl)
}
