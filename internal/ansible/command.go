package ansible

import (
	log "github.com/sirupsen/logrus"
	"os"
	"text/template"
)

type Command struct {
	Comment           string
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

var templateBash = `{{ .Comment }}
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

func (tpl *Command) GenerateCmd() {
	tmpl, _ := template.New("test").Parse(templateBash)
	err := tmpl.Execute(os.Stdout, tpl)
	if err != nil {
		log.Error(err)
	}
}
