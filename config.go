package main

import (
	"errors"
	"flag"
	"github.com/BurntSushi/toml"
	"os"
	"os/exec"
)

var configFile string

type config struct {
	AccessKeyId       string `toml:"aws_access_key_id"`
	SecretAccessKey   string `toml:"aws_secret_access_key"`
	Region            string `toml:"region"`
	RoleKey           string `toml:"role_key"`
	AnsibleCmd        string `toml:"ansible_command"`
	InventoryTemplate string `toml:"inventory_template"`
	PlaybookPath      string `toml:"playbook_path"`
}

func init() {
	flag.StringVar(&configFile, "c", "/etc/ansible-kick/ansible-kick.conf",
		"config file")
}

func loadConfig() (*config, error) {
	// Set default values
	c := new(config)
	c.RoleKey = "role"
	c.AnsibleCmd = "/usr/bin/ansible-playbook"
	c.InventoryTemplate = "/etc/ansible-kick/hosts.tmpl"
	c.PlaybookPath = "/etc/ansible/playbooks"

	_, err := toml.DecodeFile(configFile, c)
	if err != nil {
		return nil, err
	}
	if c.AccessKeyId == "" {
		return nil, errors.New("missing aws_access_key_id")
	}
	if c.SecretAccessKey == "" {
		return nil, errors.New("missing aws_secret_access_key")
	}
	if c.Region == "" {
		return nil, errors.New("missing region")
	}
	if _, err := os.Stat(c.InventoryTemplate); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New(c.InventoryTemplate +
				" inventory template file does not exist")
		}
	}

	// is the ansible cmd executable?
	_, err = exec.LookPath(c.AnsibleCmd)
	if err != nil {
		return nil, errors.New(c.AnsibleCmd + " is not executable")
	}
	return c, nil
}
