package utopia

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const customizedRepo = "customized"
const templatesDir = "config-templates"

// CustomizeDir reconfigures all repositories found in directory.
func CustomizeDir(directory string) error {

	repos, err := subDirectories(directory)
	if err != nil {
		return fmt.Errorf("failed to list repos: %v", err)
	}

	return Customize(directory, repos)
}

// Customize updates the customized repository. All repositories have to be
// located in directory, only repositories listen in repos are updated.
// It renders jinja2 templates using Ansible and creates a Makefile to apply
// the custom configuration to Kubernetes.
// Makefile targets 'make configure' and 'make deploy' are hooks to sidestep the
// default template & kubectl behavior.
func Customize(directory string, repos []string) error {

	err := copyAnsibleVars(directory, repos)
	if err != nil {
		return err
	}

	customizedPath := filepath.Join(directory, customizedRepo)

	jt := []jinja2Template{}

	for _, repo := range repos {
		if repo == customizedRepo {
			continue
		}
	}

	for _, repo := range repos {

		if repo == customizedRepo {
			continue
		}
		configTemplatesDir := filepath.Join(directory, repo, templatesDir)

		if _, err := os.Stat(configTemplatesDir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(configTemplatesDir, regenerateConfig(&jt, customizedPath, repo, directory))
		if err != nil {
			return fmt.Errorf("customization failed for repo %v: %v", repo, err)
		}

	}

	err = renderJinja2(customizedPath, jt)
	if err != nil {
		return fmt.Errorf("jinja2 rendering via ansible failed: %v", err)
	}

	for _, repo := range repos {
		dir := filepath.Join(customizedPath, repo)
		if _, err := os.Stat(filepath.Join(dir, "Makefile")); err == nil {
			err = makeConfigure(dir)
			if err != nil {
				return fmt.Errorf("make configure (%v): %v", repo, err)
			}
		}
	}

	err = generateMakefile(directory)
	if err != nil {
		return fmt.Errorf("Makefile creation failed: %v", err)
	}

	return nil
}

func subDirectories(path string) ([]string, error) {
	contents, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %v", err)
	}

	subDirectories := []string{}
	for _, content := range contents {
		if !content.IsDir() {
			continue
		}
		subDirectories = append(subDirectories, content.Name())
	}
	return subDirectories, nil
}

func regenerateConfig(jinja2Templates *[]jinja2Template, customizePath, repo, directory string) filepath.WalkFunc {
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

		if filepath.Ext(info.Name()) != jinjaSuffix {
			return copy(path, dest)
		}

		src, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to get absolute of src: %v", err)
		}
		*jinja2Templates = append(*jinja2Templates, jinja2Template{
			Src:  src,
			Dest: strings.TrimSuffix(dest, jinjaSuffix),
		})

		return nil
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
