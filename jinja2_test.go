package utopia

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestJinja(t *testing.T) {

	dest1, err := ioutil.TempFile(os.TempDir(), "jinja2_test_dest")
	if err != nil {
		t.Errorf("failed to create destination: %v", err)
	}
	defer os.Remove(dest1.Name())

	dest2, err := ioutil.TempFile(os.TempDir(), "jinja2_test_dest")
	if err != nil {
		t.Errorf("failed to create destination: %v", err)
	}
	defer os.Remove(dest2.Name())

	src, err := filepath.Abs("testdata/jinja.input")
	if err != nil {
		t.Errorf("failed to find input path: %v", err)
	}

	jt := []jinja2Template{{
		Src:  src,
		Dest: dest1.Name(),
	}, {
		Src:  src,
		Dest: dest2.Name(),
	}}
	err = renderJinja2("testdata", jt)
	if err != nil {
		t.Errorf("failed to use jinja: %v", err)
	}

	golden, err := ioutil.ReadFile("testdata/jinja.golden")
	if err != nil {
		t.Errorf("failed to read golden state: %v", err)
	}

	result1, err := ioutil.ReadFile(dest1.Name())
	if err != nil {
		t.Errorf("failed to read result: %v", err)
	}
	if bytes.Compare(result1, golden) != 0 {
		t.Errorf("Jinja was incorrect, got: %+s, want: %+s.", result1, golden)
	}

	result2, err := ioutil.ReadFile(dest2.Name())
	if err != nil {
		t.Errorf("failed to read result: %v", err)
	}
	if bytes.Compare(result2, golden) != 0 {
		t.Errorf("Jinja was incorrect, got: %+s, want: %+s.", result2, golden)
	}
}
