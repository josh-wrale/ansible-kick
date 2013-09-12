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

type AnsibleRequest struct {
	Host              string
	Role              string
	InventoryTempFile *os.File
	PlaybookFilePath  string
}

func SetRole(host string) (role string, err error) {
	auth := aws.Auth{
		config.AWS.AccessKeyId,
		config.AWS.SecretAccessKey,
	}

	region := aws.Regions[config.AWS.Region]

	e := ec2.New(auth, region)
	filter := ec2.NewFilter()
	filter.Add("private-ip-address", host)

	response, err := e.Instances(nil, filter)
	if err != nil {
		return "", errors.New("EC2 API call failed; verify AWS credentials")
	}

	if len(response.Reservations) != 1 {
		err = errors.New("no instance with private ip-address " + host)
		return "", err
	} else {
		tags := response.Reservations[0].Instances[0].Tags
		for _, tag := range tags {
			if tag.Key == config.AWS.RoleKey {
				role = tag.Value
			}
		}
	}

	instanceId := response.Reservations[0].Instances[0].InstanceId
	if role == "" {
		err = errors.New("role tag not defined on ec2 instance " + instanceId)
	}
	return role, err
}

func (ar *AnsibleRequest) SetPlaybookFilePath() error {
	fn := path.Join(config.Ansible.PlaybookPath, ar.Role+".yml")
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return errors.New("playbook does not exist " + fn)
	} else {
		ar.PlaybookFilePath = fn
	}
	return nil
}

func (ar *AnsibleRequest) SetInventoryFile() (err error) {
	ar.InventoryTempFile, err = ioutil.TempFile("", "inventory")
	if err != nil {
		return err
	}

	t, err := template.New("hosts.tmpl").ParseFiles(config.Ansible.HostsTemplate)
	if err != nil {
		return err
	}

	err = t.Execute(ar.InventoryTempFile, ar)
	if err != nil {
		return err
	}

	return
}

func (ar *AnsibleRequest) Run() (err error) {
	cmd := exec.Command(config.Ansible.Cmd, ar.PlaybookFilePath,
		"-i", ar.InventoryTempFile.Name(), "-l", ar.Host)
	out, err := cmd.CombinedOutput()
	log.Printf("%s", out)

	return err
}
