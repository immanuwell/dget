package main

import (
	"os/exec"
	"strings"
	"testing"
	"fmt"
)

type CmdOutput struct {
	Cmd    *exec.Cmd
	Output string
}

func TestCommandOutput(t *testing.T) {
	commands := []CmdOutput{
		{Cmd: exec.Command("go", "run", ".", "--debug", "example.com"),
			Output: "[example.com]"},
		{Cmd: exec.Command("go", "run", ".", "--base", "example,test", "--tld", "com,net", "--debug"),
			Output: "[example.com example.net test.com test.net]"},
		{Cmd: exec.Command("go", "run", ".", "--base", "example,test", "--tld", "com", "--tld", "net", "--debug"),
			Output: "[example.com example.net test.com test.net]"},
		{Cmd: exec.Command("go", "run", ".", "--domains-file", "domains-file.txt", "--tld", "com", "--debug"),
			Output: "[zero.com new-domain.com djkflsdl.com]"},
		{Cmd: exec.Command("go", "run", ".", "--domains-file", "domains-file.txt", "--tld-file", "tlds-file.txt", "--debug"),
			Output: "[zero.com zero.net zero.org zero.io zero.co new-domain.com new-domain.net new-domain.org new-domain.io new-domain.co djkflsdl.com djkflsdl.net djkflsdl.org djkflsdl.io djkflsdl.co]"},
	}

	for _, item := range commands {
		out, err := item.Cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("[%s] error: %v\n", item.Output, err)
			continue
		}
		want := item.Output
		output := strings.TrimSpace(string(out)) // remove trailing \n
		lines := strings.Split(output, "\n")     // split into lines
		got := lines[len(lines)-1]               // last printed line

		if got != want {
			t.Fatalf("unexpected output:\n got: %q\nwant: %q", got, want)
		}
	}

}
