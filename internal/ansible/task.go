package ansible

type Task struct {
	Role string   `yaml:"role"`
	Tags []string `yaml:"tags,omitempty" yaml:"tags,omitempty"`
}
