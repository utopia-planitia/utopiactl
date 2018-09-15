package utopia

import (
	"log"
	"path/filepath"
)

const customizedRepo = "customized"
const templatesDir = "config-templates"

func Customize(directory string, repos []string) {

	customizedPath := filepath.Join(directory, customizedRepo)

	for _, repo := range repos {

		if repo == customizedRepo {
			continue
		}

		repoPath := filepath.Join(directory, repo, templatesDir)

		err := filepath.Walk(repoPath, renderConfig(customizedPath, repo, directory))
		if err != nil {
			log.Fatalf("failed to customize %v: %v", repo, err)
		}

		// prerender for example certificates
	}

}
