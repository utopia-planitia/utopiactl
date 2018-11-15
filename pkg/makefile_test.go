package utopia

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestMakefile(t *testing.T) {

	err := generateMakefile("testdata/makefile")
	if err != nil {
		t.Errorf("failed to generate makefile: %v", err)
		return
	}

	golden, err := ioutil.ReadFile("testdata/makefile/Makefile.golden")
	if err != nil {
		t.Errorf("failed to read golden state: %v", err)
		return
	}

	result, err := ioutil.ReadFile("testdata/makefile/Makefile")
	if err != nil {
		t.Errorf("failed to read result: %v", err)
		return
	}
	if bytes.Compare(result, golden) != 0 {
		t.Errorf("Makefile was incorrect, got: %+s, want: %+s.", result, golden)
		return
	}

	defer os.Remove("testdata/makefile/Makefile")
}

func TestMakefileEmpty(t *testing.T) {

	err := generateMakefile("testdata/makefile-empty")
	if err != nil {
		t.Errorf("failed to generate makefile: %v", err)
		return
	}

	golden, err := ioutil.ReadFile("testdata/makefile-empty/Makefile.golden")
	if err != nil {
		t.Errorf("failed to read golden state: %v", err)
		return
	}

	result, err := ioutil.ReadFile("testdata/makefile-empty/Makefile")
	if err != nil {
		t.Errorf("failed to read result: %v", err)
		return
	}
	if bytes.Compare(result, golden) != 0 {
		t.Errorf("Makefile was incorrect, got: %+s, want: %+s.", result, golden)
		return
	}

	defer os.Remove("testdata/makefile-empty/Makefile")
}
