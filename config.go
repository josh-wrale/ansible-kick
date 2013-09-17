// Copyright 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.
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
	AccessKeyID       string   `toml:"aws_access_key_id"`
	SecretAccessKey   string   `toml:"aws_secret_access_key"`
	Region            string   `toml:"region"`
	RoleKey           string   `toml:"role_key"`
	AnsibleCmd        string   `toml:"ansible_command"`
	InventoryTemplate string   `toml:"inventory_template"`
	PlaybookPath      string   `toml:"playbook_path"`
	ExtraVarTags      []string `toml:"extra_var_tags"`
}

func init() {
	flag.StringVar(&configFile, "c", "/etc/ansible-kick/ansible-kick.conf",
		"config file")
}

// loadConfig sets default configuration values and loads additional settings
// from an external configuration file, overriding defaults.
// It returns a config object and an non-nil error if any of the required
// settings are missing or invalid.
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
	if c.AccessKeyID == "" {
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
