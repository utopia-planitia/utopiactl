package main

import (
	"bytes"
	"io/ioutil"
	"log"
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
	log.Println(pwd)
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
		t.Errorf("Jinja was incorrect, got: %+s, want: %+s.", result, golden)
	}
}
