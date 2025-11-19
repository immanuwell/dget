package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestCommandOutput(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "single-domain-debug",
			args: []string{"run", ".", "--debug", "example.com"},
			want: "[example.com]",
		},
		{
			name: "base-comma-tld-comma",
			args: []string{"run", ".", "--base", "example,test", "--tld", "com,net", "--debug"},
			want: "[example.com example.net test.com test.net]",
		},
		{
			name: "base-multi-tld-flags",
			args: []string{"run", ".", "--base", "example,test", "--tld", "com", "--tld", "net", "--debug"},
			want: "[example.com example.net test.com test.net]",
		},
		{
			name: "domains-file-single-tld",
			args: []string{"run", ".", "--domains-file", "domains-file.txt", "--tld", "com", "--debug"},
			want: "[zero.com new-domain.com djkflsdl.com]",
		},
		{
			name: "domains-file-and-tld-file",
			args: []string{"run", ".", "--domains-file", "domains-file.txt", "--tld-file", "tlds-file.txt", "--debug"},
			want: "[zero.com zero.net zero.org zero.io zero.co new-domain.com new-domain.net new-domain.org new-domain.io new-domain.co djkflsdl.com djkflsdl.net djkflsdl.org djkflsdl.io djkflsdl.co]",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cmd := exec.Command("go", tt.args...)

			var outBuf bytes.Buffer
			cmd.Stdout = &outBuf
			cmd.Stderr = &outBuf

			if err := cmd.Run(); err != nil {
				t.Fatalf("command %q failed: %v\noutput:\n%s", strings.Join(tt.args, " "), err, outBuf.String())
			}

			output := strings.TrimSpace(outBuf.String())
			lines := strings.Split(output, "\n")
			got := lines[len(lines)-1]

			if got != tt.want {
				t.Fatalf("unexpected output:\n got: %q\nwant: %q", got, tt.want)
			}
		})
	}
}
