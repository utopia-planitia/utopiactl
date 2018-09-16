package utopia

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
)

const customizedRepo = "customized"
const templatesDir = "config-templates"

func Customize(directory string, repos []string) error {

	customizedPath := filepath.Join(directory, customizedRepo)

	jt := []jinja2Template{}

	for _, repo := range repos {

		if repo == customizedRepo {
			continue
		}

		repoPath := filepath.Join(directory, repo, templatesDir)

		err := filepath.Walk(repoPath, parseConfig(&jt, customizedPath, repo, directory))
		if err != nil {
			return fmt.Errorf("customization failed for repo %v: %v", repo, err)
		}

		if _, err := os.Stat(filepath.Join(repoPath, "Makefile")); err == nil {
			err = makeConfigure(repoPath)
			if err != nil {
				return fmt.Errorf("make configure (%v): %v", repo, err)
			}
		}
	}

	err := renderJinja2(customizedPath, jt)
	if err != nil {
		return fmt.Errorf("jinja2 rendering via ansible failed: %v", err)
	}

	err = generateMakefile(directory, customizedPath)
	if err != nil {
		return fmt.Errorf("Makefile creation failed: %v", err)
	}

	return nil
}

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
