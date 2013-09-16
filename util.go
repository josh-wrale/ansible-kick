package main

import (
	"errors"
	"strings"
)

// extractIP returns a string representing the IP address extracted from the
// SSH_CLIENT environment variable.
func extractIP(sshClient string) (string, error) {
	fields := strings.Split(sshClient, " ")
	if len(fields) != 3 {
		return "", errors.New("invalid SSH_CLIENT format")
	}
	return fields[0], nil
}
