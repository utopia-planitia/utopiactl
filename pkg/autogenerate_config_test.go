package utopia

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAutogenerateConfigs(t *testing.T) {

	services := []string{"repo1", "repo2", "repo3", "repo4"}
	err := autogenerateConfigs("testdata/autogenerate_configs", services)
	if err != nil {
		t.Errorf("failed to autogenerate configurations: %v", err)
	}

	files := []string{
		"repo2/static",
		"repo2/variable",
		"repo4/roles/template.j2",
	}

	for _, file := range files {

		g := filepath.Join("testdata/autogenerate_configs/golden/", file)
		r := filepath.Join("testdata/autogenerate_configs/configurations/", file)

		golden, err := ioutil.ReadFile(g)
		if err != nil {
			t.Errorf("failed to read golden state: %v", err)
		}
		result, err := ioutil.ReadFile(r)
		if err != nil {
			t.Errorf("failed to read result: %v", err)
		}
		if bytes.Compare(result, golden) != 0 {
			t.Errorf("config generation was incorrect, got: %+s, want: %+s.", result, golden)
		}
	}

	if _, err := os.Stat("testdata/autogenerate_configs/configurations/repo3/touched_file"); err != nil {
		t.Errorf("failed to find file: %v", err)
	}
}
