package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	utopia "github.com/utopia-planitia/utopiactl/pkg/utopia"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to determine current working directory: %v", err)
	}

	repos, err := services(cwd, os.Args[1:])
	if err != nil {
		log.Fatalf("failed to setup config: %v", err)
	}

	err = utopia.Customize(cwd, repos)
	if err != nil {
		log.Fatalf("failed to setup config: %v", err)
	}
}

func services(directory string, args []string) ([]string, error) {
	if len(args) != 0 {
		return args, nil
	}
	services, err := subDirectories(filepath.Join(directory, "services"))
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %v", err)
	}
	return services, nil
}

func subDirectories(path string) ([]string, error) {
	contents, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %v", err)
	}

	ls := []string{}
	for _, content := range contents {
		if !content.IsDir() {
			continue
		}
		ls = append(ls, content.Name())
	}
	return ls, nil
}
