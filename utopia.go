package utopia

import (
	"fmt"
	"path/filepath"
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
			return fmt.Errorf("failed to customize %v: %v", repo, err)
		}

		// prerender for example certificates
	}

	err := renderJinja2(customizedPath, jt)
	if err != nil {
		return fmt.Errorf("failed to render jinja2 templates via ansible: %v", err)
	}

	return nil
}
