package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
)

func main() {
	//cwd := filepath.Dir(os.Args[0])
	cwd, _ := filepath.Abs("./")

	repos, err := repositories(cwd, os.Args[1:])
	if err != nil {
		log.Fatalf("failed to setup config: %v", err)
	}

	customizePath := filepath.Join(cwd, "customize")

	for _, repo := range repos {

		log.Println(repo)

		if repo == "customize" {
			continue
		}

		repoPath := filepath.Join(cwd, repo)

		err := vars(repoPath, customizePath)
		if err != nil {
			log.Fatalf("failed to copy vars from %v: %v", customizePath, err)
		}

		filepath.Walk(repoPath, walk(customizePath, repo, cwd))

		// prerender for example certificates
	}
}

func walk(customizePath, repo, cwd string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		subPath := strings.TrimPrefix(cwd, path)
		dest := filepath.Join(customizePath, repo, subPath)

		if info.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		// render jinja2 config templates

		return copy.Copy(path, dest)
	}
}
