package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigure(t *testing.T) {

	orgPWD := os.Getenv("PWD")
	defer os.Chdir(orgPWD)
	orgArgs := os.Args
	defer func() { os.Args = orgArgs }()

	pwd, err := filepath.Abs("testdata")
	if err != nil {
		t.Errorf("failed find testdata directory: %v", err)
	}
	os.Chdir(pwd)
	os.Args = []string{"utopiactl", "cfg", "all"}
	main()

	os.Chdir(orgPWD)
	os.Args = orgArgs
	os.Setenv("PWD", orgPWD)

	result, err := ioutil.ReadFile("testdata/configurations/service-repo/template")
	if err != nil {
		t.Errorf("failed to read result: %v", err)
	}
	golden, err := ioutil.ReadFile("testdata/golden/configurations/service-repo/template")
	if err != nil {
		t.Errorf("failed to read golden state: %v", err)
	}

	if bytes.Compare(result, golden) != 0 {
		t.Errorf("Jinja rendering was incorrect, got: %+s, want: %+s.", result, golden)
	}

	result, err = ioutil.ReadFile("testdata/Makefile")
	if err != nil {
		t.Errorf("failed to read result: %v", err)
	}
	golden, err = ioutil.ReadFile("testdata/golden/Makefile")
	if err != nil {
		t.Errorf("failed to read golden state: %v", err)
	}

	if bytes.Compare(result, golden) != 0 {
		t.Errorf("Makefile was incorrect, got: %+s, want: %+s.", result, golden)
	}
}

func TestExec(t *testing.T) {

	orgPWD := os.Getenv("PWD")
	defer os.Chdir(orgPWD)
	orgArgs := os.Args
	defer func() { os.Args = orgArgs }()

	os.Remove("testdata/services/service-repo2/abc")

	pwd, err := filepath.Abs("testdata")
	if err != nil {
		t.Errorf("failed find testdata directory: %v", err)
	}
	os.Chdir(pwd)
	os.Args = []string{"utopiactl", "exec", "service-repo2", "bash", "-c", "echo -n def > abc"}
	main()

	os.Chdir(orgPWD)
	os.Args = orgArgs
	os.Setenv("PWD", orgPWD)

	result, err := ioutil.ReadFile("testdata/services/service-repo2/abc")
	if err != nil {
		t.Errorf("failed to read result: %v", err)
	}

	def := []byte("def")
	if bytes.Compare(result, def) != 0 {
		t.Errorf("command execution failed, got: %+s, want: %+s.", result, def)
	}
}
