# dploy
![test](https://github.com/ca-gip/dploy/workflows/test/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/ca-gip/dploy/badge.svg)](https://coveralls.io/github/ca-gip/dploy)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ca-gip/dploy)
[![Go Report Card](https://goreportcard.com/badge/github.com/ca-gip/dploy)](https://goreportcard.com/report/github.com/ca-gip/dploy)
![dowload](https://img.shields.io/github/downloads/ca-gip/dploy/total)

`ploy` is a meta CLI for Ansible which enable multi inventory deployment while providing smart completion.

It supports `ansible-playbook` and `ansible` command

### Intallation

```bash
curl -ssf -L https://raw.githubusercontent.com/ca-gip/dploy/master/install.sh | bash
```

Adding completion for bash (also available for zsh and fish)

```
echo "source <(dploy completion bash)" >> ~/.bashrc
```

## Usage

```bash
# dploy -h
Ansible deployment toolbox

Usage:
  dploy [command]

Available Commands:
  completion  Generate completion script
  exec        Run Ad Hoc command
  generate    Generate ansible-playbook command
  help        Help about any command
  play        Run ansible-playbook command
```

### How to select inventories ?

All subcommand use a `--filter` arguments that will select inventories based on the vars declared under `[all:vars]` in *.ini files.

Filter implement the following operators to match variable value:
 * Equal : `==`
 * NotEqual `!=`
 * EndWith `$=`
 * Contains `~=`
 * StartWith `^=`

Filtering is based on the location where the command is executed as it will recursively search all *.ini files.

![filtering](https://asciinema.org/a/CJe2VHw7eftbp5J1onjBSsdtT.svg)

### Play subcommand

Launch a playbook
```bash
dploy play --filter platform==os -p upgrade-proxy.yml
```

### Exec subcommand

Execute Ad Hoc command 
```bash
dploy exec --filter customer==os -m ping -p lb
```