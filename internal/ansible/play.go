package ansible

import (
	"github.com/ca-gip/dploy/internal/utils"
)

type Play struct {
	Hosts string    `yaml:"hosts"`
	Roles []*Role   `yaml:"roles"`
	Tags  utils.Set `yaml:"tags"`
}

func (play *Play) AllTags() (tags *utils.Set) {
	tags = utils.NewSet()
	for _, role := range play.Roles {
		tags.Concat(role.AllTags().List())
	}
	tags.Concat(play.Tags.List())
	return
}
