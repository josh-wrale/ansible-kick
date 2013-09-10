package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	AWS     AWSConfig
	Ansible AnsibleConfig
}

type AWSConfig struct {
	AccessKeyId     string `toml:"access_key_id"`
	SecretAccessKey string `toml:"secret_access_key"`
	Region          string `toml:"region"`
	RoleKey         string `toml:"role_key"`
}

type AnsibleConfig struct {
	Cmd           string `toml:"command"`
	HostsTemplate string `toml:"hosts_template"`
	PlaybookPath  string `toml:"playbook_path"`
}

func DecodeFile(fpath string, v interface{}) error {
	_, err := toml.DecodeFile(c, v)
	return err
}
