package ansible

type Play struct {
	Hosts string   `yaml:"hosts"`
	Roles []Role   `yaml:"roles"`
	Tags  []string `yaml:"tags,omitempty"`
}
