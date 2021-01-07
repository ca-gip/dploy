package services

type Task struct {
	Role string      `yaml:"role"`
	Tags interface{} `yaml:"tags,omitempty" yaml:"tags,omitempty"`
}
