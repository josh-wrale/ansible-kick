// Copyright 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.
package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"text/template"
)

type ansibleConfig struct {
	Host              string
	Role              string
	inventoryTempFile *os.File
	playbookFilePath  string
	settings          *config
}

// setRole sets the Ansible role based on ec2 tags of the target instance.
// If the target ec2 instance does not exist, or does not have a role tag
// value, then setRole returns an error explaining why the role could not
// be set.
func (ac *ansibleConfig) setRole() error {
	i, err := ac.findInstance()
	if err != nil {
		return err
	}

	for _, tag := range i.Tags {
		if tag.Key == ac.settings.RoleKey {
			ac.Role = tag.Value
		}
	}
	if ac.Role == "" {
		return errors.New("role tag not defined on ec2 instance " + i.InstanceId)
	}
	return nil
}

// setPlaybookFilePath sets the path of the Ansible playbook file and
// returns an error if the playbook does not exist.
func (ac *ansibleConfig) setPlaybookFilePath() error {
	fn := path.Join(ac.settings.PlaybookPath, ac.Role+".yml")
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return errors.New("playbook does not exist " + fn)
	} else {
		ac.playbookFilePath = fn
	}
	return nil
}

// setInventoryFile sets the path of the Ansible inventory file to a
// dynamically generated inventory file using an inventory template.
// It returns an error if the inventory file cannot be generated.
func (ac *ansibleConfig) setInventoryFile() error {
	var err error
	ac.inventoryTempFile, err = ioutil.TempFile("", "inventory")
	if err != nil {
		return err
	}
	t, err := template.New("hosts.tmpl").ParseFiles(ac.settings.InventoryTemplate)
	if err != nil {
		return err
	}
	if err := t.Execute(ac.inventoryTempFile, ac); err != nil {
		return err
	}
	return nil
}

// run executes an ansible-playbook run.
// It returns any error generated by the ansible-playbook command.
func (ac *ansibleConfig) run() error {
	cmd := exec.Command(ac.settings.AnsibleCmd, ac.playbookFilePath,
		"-i", ac.inventoryTempFile.Name(), "-l", ac.Host)
	out, err := cmd.CombinedOutput()
	log.Printf("%s", out)
	return err
}

// Run prepares and executes an ansible-playbook run for a specific host.
// It returns an error if any part of the preparation or execution fails.
func Run(host string) error {
	var err error
	ac := new(ansibleConfig)
	ac.Host = host
	ac.settings, err = loadConfig()
	if err != nil {
		return err
	}
	if err := ac.setRole(); err != nil {
		return err
	}

	ac.setInventoryFile()
	defer os.Remove(ac.inventoryTempFile.Name())
	defer ac.inventoryTempFile.Close()

	if err := ac.setPlaybookFilePath(); err != nil {
		return err
	}
	if err := ac.run(); err != nil {
		return err
	}
	return nil
}
