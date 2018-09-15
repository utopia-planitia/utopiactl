package utopia

import (
	"log"
	"path/filepath"
)

const customizedRepo = "customized"

func Customize(directory string, repos []string) {

	customizePath := filepath.Join(directory, customizedRepo)

	for _, repo := range repos {

		log.Println(repo)

		if repo == "customize" {
			continue
		}

		repoPath := filepath.Join(directory, repo)

		err := filepath.Walk(repoPath, renderConfig(customizePath, repo, directory))
		if err != nil {
			log.Fatalf("failed to customize %v: %v", repo, err)
		}

		// prerender for example certificates
	}

}
