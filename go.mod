module github.com/ca-gip/dploy

go 1.15

replace github.com/apenella/go-ansible => github.com/clementblaise/go-ansible v0.6.2-0.20210121132918-f754b400712f

require (
	github.com/apenella/go-ansible v0.6.1
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-test/deep v1.0.7
	github.com/karrick/godirwalk v1.16.1
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pterm/pterm v0.12.8
	github.com/relex/aini v1.2.1
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897
	gopkg.in/yaml.v2 v2.3.0
)
