package ansible

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"text/template"
)

type PlaybookCmd struct {
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

var PlaybookCmdTemplate = `{{ .Comment }}
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

func (tpl *PlaybookCmd) Generate() {
	tmpl, _ := template.New("test").Parse(PlaybookCmdTemplate)
	err := tmpl.Execute(os.Stdout, tpl)
	if err != nil {
		log.Error(err)
	}
}

type AdHocCmd struct {
	Comment           string
	Inventory         []*Inventory
	Pattern           string
	ModuleName        string
	ModuleArgs        string
	ExtraVars         map[string]interface{}
	Background        int
	Fork              int
	PollInterval      int
	Limit             []string
	Check             bool
	Diff              bool
	OneLine           bool
	Tree              bool
	PlaybookDir       string
	VaultPasswordFile string
	AskVaultPass      bool
}

var AdHocCmdTemplate = `{{ .Comment }}
{{- range $inventory := .Inventory }}
ansible {{ $.Pattern }} -a {{ $.ModuleArgs }} 
{{- if $.ModuleName }} -m {{ $.ModuleName }}{{- end }}
-i {{ $inventory.RelativePath }}
{{- if $.ExtraVars }} -e {{ $.ExtraVars }}{{- end }}
{{- if $.Background }} --background {{ $.Background }}{{- end }}
{{- if $.Fork }} --forks {{ $.Fork }}{{- end }}
{{- if $.PollInterval }} --poll {{ $.PollInterval }}{{- end }}
{{- if $.Limit }} -l {{ range $i,$limit := $.Limit }}{{ if gt $i 0 }},{{ end }}{{ $limit }}{{ end }}{{- end }}
{{- if $.Check }} -c{{- end }}
{{- if $.Diff }} --diff{{- end }}
{{- if $.OneLine }} --one-line{{- end }}
{{- if $.Tree }} --tree{{- end }}
{{- if $.PlaybookDir }} --playbook-dir {{ $.PlaybookDir }}{{- end }}
{{- if $.AskVaultPass }} --ask-vault-password{{- end }}
{{- if $.VaultPasswordFile }} --vault-password-file {{ $.VaultPasswordFile }}{{- end }}
{{- end }}
`

func (o *AdHocCmd) Generate() {
	tmpl, _ := template.New("test").Parse(AdHocCmdTemplate)
	err := tmpl.Execute(os.Stdout, o)
	if err != nil {
		log.Error(err)
	}
}

func (o *AdHocCmd) AddExtraVar(name string, value interface{}) error {

	if o.ExtraVars == nil {
		o.ExtraVars = map[string]interface{}{}
	}
	_, exists := o.ExtraVars[name]
	if exists {
		return errors.New(fmt.Sprintf("ExtraVar '%s' already exist", name))
	}

	o.ExtraVars[name] = value

	return nil
}
