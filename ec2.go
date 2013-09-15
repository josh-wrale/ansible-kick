// Copyright 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.
package main

import (
	"errors"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/ec2"
)

func searchByFilter(e *ec2.EC2, key, value string) (*ec2.Instance, error) {
	filter := ec2.NewFilter()
	filter.Add(key, value)
	resp, err := e.Instances(nil, filter)
	if err != nil {
		return nil, errors.New("EC2 API call failed. " +
			"Check your AWS credentials and system clock")
	}
	if len(resp.Reservations) != 1 {
		return nil, errors.New("no instance with " + key + "=" + value)
	} else {
		return &resp.Reservations[0].Instances[0], nil
	}
}

func (ac *ansibleConfig) findEC2Instance() (*ec2.Instance, error) {
	var (
		err      error
		instance *ec2.Instance
	)

	auth := aws.Auth{
		ac.settings.AccessKeyId,
		ac.settings.SecretAccessKey,
	}
	region := aws.Regions[ac.settings.Region]
	e := ec2.New(auth, region)

	instance, err = searchByFilter(e, "private-ip-address", ac.Host)
	if err != nil {
		instance, err = searchByFilter(e, "ip-address", ac.Host)
	}

	if err != nil {
		return nil, err
	} else {
		return instance, nil
	}
}
