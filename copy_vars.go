package utopia

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func copyVars(directory, customizedPath, repo, vars string) error {

	source := filepath.Join(directory, repo, vars)
	target := filepath.Join(customizedPath, vars)
	if _, err := os.Stat(source); err == nil {
		err := filepath.Walk(source, regenerateVars(source, target))
		if err != nil {
			return fmt.Errorf("%s sync failed for repo %v: %v", vars, repo, err)
		}
	}

	return nil
}

func regenerateVars(source, target string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
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
}
