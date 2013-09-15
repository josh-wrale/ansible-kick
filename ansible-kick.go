package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

func getHost() string {
	host := os.Getenv("SSH_CLIENT")
	host = strings.Split(host, " ")[0]
	if host == "" {
		log.Fatal("ipaddress required")
	}
	return host
}

func main() {
	flag.Parse()
	host := getHost()
	log.SetFlags(0)
	log.Print("starting ansible-kick for " + host)
	if err := Run(host); err != nil {
		log.Fatal(err.Error())
	}
	log.Print("successfully kicked " + host)
}
