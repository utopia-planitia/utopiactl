package utopia

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyAnsibleVars(t *testing.T) {

	err := removeContents("testdata/ansible_vars/ansible")
	if err != nil {
		t.Errorf("failed to cleanup: %v", err)
	}

	services := []string{"repo1", "repo2", "repo3"}
	err = copyAnsibleVars("testdata/ansible_vars", services)
	if err != nil {
		t.Errorf("failed to copy ansible vars: %v", err)
	}

	files := []string{
		"testdata/ansible_vars/ansible/host_vars/host2.yaml",
		"testdata/ansible_vars/ansible/group_vars/group3.yaml",
	}

	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			t.Errorf("failed to find file: %v", err)
		}
	}
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
