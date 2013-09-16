// Copyright 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.
package main

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	flag.Set("c", "does-not-exist")
	flag.Parse()
	expect := "open does-not-exist: no such file or directory"
	_, err := loadConfig()
	if err.Error() != expect {
		t.Errorf("expected %s, go %s", expect, err.Error())
	}

	tempConfig, _ := ioutil.TempFile("", "config_test")
	defer os.Remove(tempConfig.Name())
	defer tempConfig.Close()

	flag.Set("c", tempConfig.Name())
	flag.Parse()
	expect = "missing aws_access_key_id"
	_, err = loadConfig()
	if err.Error() != expect {
		t.Errorf("expected %s, go %s", expect, err.Error())
	}
}
