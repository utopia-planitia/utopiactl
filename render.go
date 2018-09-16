package utopia

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
)

func renderConfig(customizePath, repo, directory string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		prefix := filepath.Join(directory, repo, templatesDir)
		subPath := strings.TrimPrefix(path, prefix)
		dest := filepath.Join(customizePath, repo, subPath)

		if info.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		if filepath.Ext(info.Name()) == jinjaSuffix {
			return renderJinja2(customizePath, path, dest)
		}

		return copy.Copy(path, dest)
	}
}
