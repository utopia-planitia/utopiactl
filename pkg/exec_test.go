package utopia

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestExec(t *testing.T) {

	os.Remove("testdata/exec/services/service1/abc")

	Exec("testdata/exec", []string{"service1"}, []string{"bash", "-c", "echo -n def > abc"})

	result, err := ioutil.ReadFile("testdata/exec/services/service1/abc")
	if err != nil {
		t.Errorf("failed to read result: %v", err)
	}

	def := []byte("def")
	if bytes.Compare(result, def) != 0 {
		t.Errorf("command execution failed, got: %+s, want: %+s.", result, def)
	}
}
