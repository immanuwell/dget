package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCommandOutput(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "--base", "example,test", "--tld", "com,net", "--debug")

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\noutput:\n%s", err, out)
	}

	output := strings.TrimSpace(string(out)) // remove trailing \n
	lines := strings.Split(output, "\n")     // split into lines
	got := lines[len(lines)-1]               // last printed line
	want := "[example.com example.net test.com test.net]"

	if got != want {
		t.Fatalf("unexpected output:\n got: %q\nwant: %q", got, want)
	}
}
