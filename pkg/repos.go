package utopia

import (
	"fmt"
	"io/ioutil"
)

func Repositories(cwd string, args []string) ([]string, error) {
	if len(args) != 0 {
		return args, nil
	}
	repos, err := subDirectories(cwd)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %v", err)
	}
	return repos, nil
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
