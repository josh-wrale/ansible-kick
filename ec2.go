// Copyright 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.
package main

import (
	"errors"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/ec2"
)

type EC2Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

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
	}
	return &resp.Reservations[0].Instances[0], nil
}

func (c *EC2Config) findInstance(host string) (*ec2.Instance, error) {
	var (
		err         error
		ec2Instance *ec2.Instance
	)

	auth := aws.Auth{c.AccessKeyID, c.SecretAccessKey}
	region := aws.Regions[c.Region]
	e := ec2.New(auth, region)

	ec2Instance, err = searchByFilter(e, "private-ip-address", host)
	if err != nil {
		ec2Instance, err = searchByFilter(e, "ip-address", host)
	}
	if err != nil {
		return nil, err
	}
	return ec2Instance, nil
}
