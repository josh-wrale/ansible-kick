package main

import "testing"

var (
	sshClientString = "203.0.113.100 4321 22"
	invalidSSHClientString = "203.0.113.100"
)

func TestExtractIP(t *testing.T) {
	expect := "203.0.113.100"
	got, _ := extractIP(sshClientString)
	if got != expect {
		t.Errorf("expected %s, got %s", expect, got)
	}
	got, err := extractIP(invalidSSHClientString)
	if err.Error() != "invalid SSH_CLIENT format" {
		t.Error("Expected err invalid SSH_CLIENT format, got ", err.Error())
	}
	if got != "" {
		t.Errorf("expected empty string, got %s", got)
	}
}
