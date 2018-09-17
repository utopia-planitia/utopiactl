package utopia

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const customizedRepo = "customized"
const templatesDir = "config-templates"

// Customize updates the customized repository. All repositories are located in
// directory, only repositories listen in repos are updated. It renders jinja2
// templates using Ansible and creates a Makefile to apply the custom
// configuration to Kubernetes.
// Makefile targets 'make configure' and 'make deploy' are hooks to sidestep the
// default template & kubectl behavior.
func Customize(directory string, repos []string) error {

	customizedPath := filepath.Join(directory, customizedRepo)

	jt := []jinja2Template{}

	for _, repo := range repos {

		if repo == customizedRepo {
			continue
		}

		configTemplatesDir := filepath.Join(directory, repo, templatesDir)
		if _, err := os.Stat(configTemplatesDir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(configTemplatesDir, parseConfig(&jt, customizedPath, repo, directory))
		if err != nil {
			return fmt.Errorf("customization failed for repo %v: %v", repo, err)
		}

		if _, err := os.Stat(filepath.Join(configTemplatesDir, "Makefile")); err == nil {
			err = makeConfigure(configTemplatesDir)
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

		return copy(path, dest)
	}
}

func copy(src, dest string) error {

	info, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("could not stat %s: %s", src, err)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create file %s: %s", dest, err)
	}
	defer f.Close()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return fmt.Errorf("could not set permissions for file %s: %s", f.Name(), err)
	}

	s, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file %s: %s", src, err)
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	if err != nil {
		return fmt.Errorf("failed to copy from source to destination: %s", err)
	}

	return nil
}
