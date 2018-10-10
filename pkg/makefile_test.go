package utopia

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestMakefile(t *testing.T) {

	generateMakefile("testdata/makefile")

	golden, err := ioutil.ReadFile("testdata/makefile/Makefile.golden")
	if err != nil {
		t.Errorf("failed to read golden state: %v", err)
	}

	result, err := ioutil.ReadFile("testdata/makefile/Makefile")
	if err != nil {
		t.Errorf("failed to read result: %v", err)
	}
	if bytes.Compare(result, golden) != 0 {
		t.Errorf("Makefile was incorrect, got: %+s, want: %+s.", result, golden)
	}
}
