package main

import (
	"github.com/BurntSushi/toml"
	"log"
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

func (c *Config) Load(fpath string) error {
	// Set default values
	c.AWS.RoleKey = "role"
	c.Ansible.Cmd = "/usr/bin/ansible"
	c.Ansible.HostsTemplate = "/etc/ansible-kick/hosts.tmpl"
	c.Ansible.PlaybookPath = "/etc/ansible/playbooks"

	_, err := toml.DecodeFile(fpath, c)

	if c.AWS.AccessKeyId == "" {
		log.Fatal("missing AWS access_key_id")
	}

	if c.AWS.SecretAccessKey == "" {
		log.Fatal("missing AWS secret_access_key")
	}

	if c.AWS.Region == "" {
		log.Fatal("missing AWS region")
	}

	return err
}
