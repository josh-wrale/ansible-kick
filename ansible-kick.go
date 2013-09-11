package main

import (
	"flag"
	"log"
	"log/syslog"
	"os"
	"strings"
)

var (
	c      string
	config Config
)

func init() {
	flag.StringVar(&c, "c", "/etc/ansible-kick/ansible-kick.conf", "config file")
}

func main() {
	flag.Parse()

	logger, err := syslog.New(syslog.LOG_INFO, "ansible-kick")
	if err != nil {
		log.Fatal(err.Error())
	}

	host := os.Getenv("SSH_CLIENT")
	if host == "" {
		log.Fatal("host ip address required.")
	}

	host = strings.Split(host, " ")[0]

	err = logger.Notice("starting ansible-kick for " + host)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = config.Load(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	role, err := SetRole(host)
	if err != nil {
		log.Fatal(err.Error())
	}

	logger.Notice(role + " role selected for " + host)

	ar := AnsibleRequest{
		Host: host,
		Role: role,
	}

	ar.SetInventoryFile()
	defer os.Remove(ar.InventoryTempFile.Name())
	defer ar.InventoryTempFile.Close()

	err = ar.SetPlaybookFilePath()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = ar.Run()
	if err != nil {
		log.Fatal(err.Error())
	}

	logger.Notice("successfully kicked " + host)
}
