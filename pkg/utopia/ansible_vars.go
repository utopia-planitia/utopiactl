package utopia

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func copyAnsibleVars(directory string, services []string) error {

	dest := filepath.Join(directory, "ansible", "host_vars")
	for _, svc := range services {
		src := filepath.Join(directory, "services", svc, "host_vars")
		if _, err := os.Stat(src); err == nil {
			err := mergeCopy(src, dest)
			if err != nil {
				return fmt.Errorf("failed to copy %s: %s", src, err)
			}
		}
	}

	dest = filepath.Join(directory, "ansible", "group_vars")
	for _, svc := range services {
		src := filepath.Join(directory, "services", svc, "group_vars")
		if _, err := os.Stat(src); err == nil {
			err := mergeCopy(src, dest)
			if err != nil {
				return fmt.Errorf("failed to copy %s: %s", src, err)
			}
		}
	}

	return nil
}

func mergeCopy(source, target string) error {

	cp := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dest := filepath.Join(target, strings.TrimPrefix(path, source))

		if info.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		return copy(path, dest)
	}

	err := filepath.Walk(source, cp)
	if err != nil {
		return fmt.Errorf("could walk through filepath %s: %v", source, err)
	}

	return nil
}
