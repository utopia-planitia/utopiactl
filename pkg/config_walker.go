package utopia

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
)

const jinjaSuffix = ".j2"

func Walk(customizePath, repo, cwd string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		subPath := strings.TrimPrefix(cwd, path)
		dest := filepath.Join(customizePath, repo, subPath)

		if info.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		if strings.HasSuffix(info.Name(), jinjaSuffix) {
			renderJinja2(customizePath, path, dest)
		}

		return copy.Copy(path, dest)
	}
}
