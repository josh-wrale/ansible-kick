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

func (ac *ansibleConfig) setplaybookFilePath() error {
	fn := path.Join(ac.settings.PlaybookPath, ac.role+".yml")
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return errors.New("playbook does not exist " + fn)
	} else {
		ac.playbookFilePath = fn
	}
	return nil
}

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

func (ac *ansibleConfig) run() error {
	cmd := exec.Command(ac.settings.AnsibleCmd, ac.playbookFilePath,
		"-i", ac.inventoryTempFile.Name(), "-l", ac.host)
	out, err := cmd.CombinedOutput()
	log.Printf("%s", out)
	return err
}

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

	if err := ac.setplaybookFilePath(); err != nil {
		return err
	}
	if err := ac.run(); err != nil {
		return err
	}
	return nil
}
