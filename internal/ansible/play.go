package ansible

import (
	"github.com/ca-gip/dploy/internal/utils"
)

type Play struct {
	Hosts   string      `yaml:"hosts"`
	Roles   []Role      `yaml:"roles"`
	RawTags interface{} `yaml:"tags,omitempty"`
}

func (play *Play) Tags() []string {
	return utils.InferSlice(play.RawTags)
}
