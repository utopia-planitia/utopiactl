package utopia

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
)

func parseConfig(jinja2Templates *[]jinja2Template, customizePath, repo, directory string) filepath.WalkFunc {
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
			src, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("failed to get absolute of src: %v", err)
			}
			dest := strings.TrimSuffix(dest, jinjaSuffix)
			*jinja2Templates = append(*jinja2Templates, jinja2Template{
				Src:  src,
				Dest: dest,
			})
		}

		return copy.Copy(path, dest)
	}
}
