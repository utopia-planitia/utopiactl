package utopia

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestExec(t *testing.T) {

	os.Remove("testdata/exec/services/service1/abc")

	err := Exec("testdata/exec", []string{"service1"}, []string{"bash", "-c", "echo -n def > abc"})
	if err != nil {
		t.Errorf("failed to exec: %v", err)
	}

	result, err := ioutil.ReadFile("testdata/exec/services/service1/abc")
	if err != nil {
		t.Errorf("failed to read result: %v", err)
	}

	def := []byte("def")
	if !bytes.Equal(result, def) {
		t.Errorf("command execution failed, got: %+s, want: %+s.", result, def)
	}
}
