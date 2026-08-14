package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteCompletionZshRegistersFunction(t *testing.T) {
	var out bytes.Buffer
	if err := writeCompletion([]string{"zsh"}, "/unused", &out); err != nil {
		t.Fatalf("writeCompletion: %v", err)
	}

	scriptPath := filepath.Join(t.TempDir(), "completion.zsh")
	if err := os.WriteFile(scriptPath, out.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(
		"zsh", "-f", "-c",
		`source "$1"; whence -w _ssmctl; print -r -- $_comps[ssmctl]`,
		"ssmctl-completion-test", scriptPath,
	)
	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("source completion script: %v\n%s", err, got)
	}
	if want := "_ssmctl: function\n_ssmctl\n"; string(got) != want {
		t.Fatalf("zsh registration output = %q, want %q", got, want)
	}
}

func TestWriteCompletionZshUsesSubstringMatcherForInstances(t *testing.T) {
	var script bytes.Buffer
	if err := writeCompletion([]string{"zsh"}, "/unused", &script); err != nil {
		t.Fatalf("writeCompletion: %v", err)
	}

	tempDir := t.TempDir()
	scriptPath := filepath.Join(tempDir, "completion.zsh")
	if err := os.WriteFile(scriptPath, script.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}

	binDir := filepath.Join(tempDir, "bin")
	if err := os.Mkdir(binDir, 0o700); err != nil {
		t.Fatal(err)
	}
	fakeSSMCTL := filepath.Join(binDir, "ssmctl")
	if err := os.WriteFile(
		fakeSSMCTL,
		[]byte("#!/bin/sh\nprintf 'live-lg-api-01\\nother-host\\n'\n"),
		0o700,
	); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(
		"zsh", "-f", "-c",
		`source "$1"
_describe() { return 0 }
compadd() { print -rl -- "$@" }
words=(ssmctl api)
CURRENT=2
_ssmctl`,
		"ssmctl-substring-test", scriptPath,
	)
	cmd.Env = append(os.Environ(), "PATH="+binDir+":"+os.Getenv("PATH"))
	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exercise completion script: %v\n%s", err, got)
	}
	if want := "-M\nl:|=* r:|=*\n-a\ninstances\n"; string(got) != want {
		t.Fatalf("compadd arguments = %q, want %q", got, want)
	}
}

func TestWriteCompletionInstancesUsesConfigAndSorts(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yml")
	data := []byte("instances:\n  z-app:\n    target: i-z\n  a-app:\n    target: i-a\n")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := writeCompletion([]string{"__instances"}, path, &out); err != nil {
		t.Fatalf("writeCompletion: %v", err)
	}
	if got, want := out.String(), "a-app\nz-app\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestWriteCompletionInstancesSilencesConfigErrors(t *testing.T) {
	invalidPath := filepath.Join(t.TempDir(), "invalid.yml")
	if err := os.WriteFile(invalidPath, []byte("instances: ["), 0o600); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{"/missing/config.yml", invalidPath} {
		t.Run(path, func(t *testing.T) {
			var out bytes.Buffer
			if err := writeCompletion([]string{"__instances"}, path, &out); err != nil {
				t.Fatalf("writeCompletion: %v", err)
			}
			if out.Len() != 0 {
				t.Fatalf("output = %q, want empty", out.String())
			}
		})
	}
}

func TestWriteCompletionRejectsInvalidArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "missing shell"},
		{name: "unsupported shell", args: []string{"bash"}},
		{name: "extra argument", args: []string{"zsh", "extra"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			err := writeCompletion(tt.args, "/unused", &out)
			if err == nil || !strings.Contains(err.Error(), "usage: ssmctl completion <zsh>") {
				t.Fatalf("error = %v, want completion usage", err)
			}
		})
	}
}
