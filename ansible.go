// Copyright 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.
package main

import (
	"errors"
	"io/ioutil"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/ec2"
	"log"
	"os"
	"os/exec"
	"path"
	"text/template"
)

type ansibleConfig struct {
	host              string
	role              string
	inventoryTempFile *os.File
	playbookFilePath  string
	settings          *config
}

// setRole sets the Ansible role based on ec2 tags of the target instance.
// If the target ec2 instance does not exist, or does not have a role tag
// value, then setRole returns an error explaining why the role could not
// be set.
func (ac *ansibleConfig) setRole() error {
	auth := aws.Auth{
		ac.settings.AccessKeyId,
		ac.settings.SecretAccessKey,
	}
	region := aws.Regions[ac.settings.Region]
	e := ec2.New(auth, region)
	filter := ec2.NewFilter()
	filter.Add("private-ip-address", ac.host)

	response, err := e.Instances(nil, filter)
	if err != nil {
		return errors.New("EC2 API call failed. " +
			"Check your AWS credentials and system clock")
	}
	if len(response.Reservations) != 1 {
		return errors.New("no instance with private ip-address " +
			ac.host)
	} else {
		tags := response.Reservations[0].Instances[0].Tags
		for _, tag := range tags {
			if tag.Key == ac.settings.RoleKey {
				ac.role = tag.Value
			}
		}
	}
	instanceId := response.Reservations[0].Instances[0].InstanceId
	if ac.role == "" {
		err = errors.New("role tag not defined on ec2 instance " + instanceId)
	}
	return err
}

// setPlaybookFilePath sets the path of the Ansible playbook file and
// returns an error if the playbook does not exist.
func (ac *ansibleConfig) setPlaybookFilePath() error {
	fn := path.Join(ac.settings.PlaybookPath, ac.role+".yml")
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
		"-i", ac.inventoryTempFile.Name(), "-l", ac.host)
	out, err := cmd.CombinedOutput()
	log.Printf("%s", out)
	return err
}

// Run prepares and executes an ansible-playbook run for a specific host.
// It returns an error if any part of the preparation or execution fails.
func Run(host string) error {
	var err error
	ac := new(ansibleConfig)
	ac.host = host
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
