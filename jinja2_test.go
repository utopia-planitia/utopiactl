package utopia

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestJinja(t *testing.T) {

	dest, err := ioutil.TempFile(os.TempDir(), "jinja2_test_dest")
	if err != nil {
		t.Errorf("failed to create destination: %v", err)
	}
	defer os.Remove(dest.Name())

	err = renderJinja2("testdata", "testdata/jinja.input", dest.Name())
	if err != nil {
		t.Errorf("failed to use jinja: %v", err)
	}

	result, err := ioutil.ReadFile(dest.Name())
	if err != nil {
		t.Errorf("failed to read result: %v", err)
	}
	golden, err := ioutil.ReadFile("testdata/jinja.golden")
	if err != nil {
		t.Errorf("failed to read golden state: %v", err)
	}

	if bytes.Compare(result, golden) != 0 {
		t.Errorf("Jinja was incorrect, got: %+s, want: %+s.", result, golden)
	}
}
