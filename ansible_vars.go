package utopia

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func copyAnsibleVars(repoPath, customizedPath, vars string) error {

	source := filepath.Join(repoPath, vars)
	target := filepath.Join(customizedPath, vars)

	if _, err := os.Stat(source); err != nil {
		return nil
	}

	cp := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		subPath := strings.TrimPrefix(path, source)
		dest := filepath.Join(target, subPath)

		if info.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		return copy(path, dest)
	}

	if _, err := os.Stat(source); err == nil {
		err := filepath.Walk(source, cp)
		if err != nil {
			return fmt.Errorf("%s sync failed for repo %v: %v", vars, repo, err)
		}
	}

	return nil
}
