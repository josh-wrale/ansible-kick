// Copyright 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.
package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	flag.Parse()

	// Don't trust the client to send it's own IP address; instead extract
	// the IP from the SSH_CLIENT environment variable, which contains three
	// space-separated pieces of information.  The IP address and source
	// port of the client followed by the server's destination port number.
	//
	//  SSH_CLIENT = 203.0.113.10 4532 22
	host, err := extractIP(os.Getenv("SSH_CLIENT"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Disable printing of the timestamp and hostname when logging to the
	// console.
	log.SetFlags(0)

	log.Print("starting ansible-kick for " + host)
	if err := Run(host); err != nil {
		log.Fatal(err.Error())
	}
	log.Print("successfully kicked " + host)
}
