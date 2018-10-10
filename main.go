package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/utopia-planitia/utopiactl/pkg"
)

const help = `usage:
	utopiactl configure [service-selector]
	utopiactl exec [service-selector] [command]

how to add a service:
	git submodule add git@gitlab.com:utopia-planitia/kured.git services/kured
	utopiactl configure kured
	git commit -a -m "added kured (kubernetes reboot daemon)"
	git push origin master

how to update a service:
	utopiactl exec kured git pull
	utopiactl configure kured
	git commit -a -m "updated kured (kubernetes reboot daemon)"
	git push origin master
`

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to determine current working directory: %v", err)
	}

	if len(os.Args) < 3 {
		printHelp()
		return
	}

	command := os.Args[1]

	svcs, err := services(cwd, os.Args[2])
	if err != nil {
		log.Fatalf("failed to select services: %v", err)
	}

	if contains([]string{"configure", "reconfigure", "config", "conf", "cfg", "c"}, command) {
		err := utopia.Configure(cwd, svcs)
		if err != nil {
			log.Fatalf("failed to auto configure: %v", err)
		}
		return
	}

	if len(os.Args) < 4 {
		printHelp()
		return
	}

	if contains([]string{"execute", "exec", "exe", "e"}, command) {
		err := utopia.Exec(cwd, svcs, os.Args[3:])
		if err != nil {
			log.Fatalf("failed to execute: %v", err)
		}
		return
	}

	printHelp()
}

func printHelp() {
	os.Stdout.WriteString(help)
}

func services(directory string, ls string) ([]string, error) {
	if ls != "" && ls != "all" {
		return strings.Split(ls, ","), nil
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
