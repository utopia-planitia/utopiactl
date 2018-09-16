package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestUtopia(t *testing.T) {

	os.RemoveAll("testdata/input/customized")

	orgPWD := os.Getenv("PWD")
	defer os.Chdir(orgPWD)
	osArgs := os.Args
	defer func() { os.Args = osArgs }()

	pwd, err := filepath.Abs("testdata/input")
	if err != nil {
		t.Errorf("failed find testdata directory: %v", err)
	}
	os.Chdir(pwd)
	os.Args = []string{"utopia"}
	main()

	result, err := ioutil.ReadFile("customized/service-repo/template")
	if err != nil {
		t.Errorf("failed to read result: %v", err)
	}
	golden, err := ioutil.ReadFile("../golden/customized/service-repo/template")
	if err != nil {
		t.Errorf("failed to read golden state: %v", err)
	}

	if bytes.Compare(result, golden) != 0 {
		t.Errorf("Jinja rendering was incorrect, got: %+s, want: %+s.", result, golden)
	}

	result, err = ioutil.ReadFile("customized/Makefile")
	if err != nil {
		t.Errorf("failed to read result: %v", err)
	}
	golden, err = ioutil.ReadFile("../golden/customized/Makefile")
	if err != nil {
		t.Errorf("failed to read golden state: %v", err)
	}

	if bytes.Compare(result, golden) != 0 {
		t.Errorf("Makefile was incorrect, got: %+s, want: %+s.", result, golden)
	}
}
