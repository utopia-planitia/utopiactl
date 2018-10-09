package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestUtopia(t *testing.T) {

	orgPWD := os.Getenv("PWD")
	defer os.Chdir(orgPWD)
	orgArgs := os.Args
	defer func() { os.Args = orgArgs }()

	pwd, err := filepath.Abs("testdata/input")
	if err != nil {
		t.Errorf("failed find testdata directory: %v", err)
	}
	os.Chdir(pwd)
	os.Args = []string{"utopiactl"}
	main()

	os.Chdir(orgPWD)
	os.Args = orgArgs
	os.Setenv("PWD", orgPWD)

	result, err := ioutil.ReadFile("testdata/input/configurations/service-repo/template")
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

	result, err = ioutil.ReadFile("testdata/input/Makefile")
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
