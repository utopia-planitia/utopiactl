package main

import (
	"log"
	"os"

	utopia "github.com/utopia-planitia/utopiactl/pkg"
)

const help = `usage:
	utopiactl configure [service-selector]
	utopiactl exec [service-selector] [command]
	utopiactl deploy [service-selector]
	utopiactl compare-versions [service-selector]

service-selector:
	kubed: selects "kubed"
	kubed,logging,metrics: selects any service listed
	all: selects all service folders found
	-: select the cluster itself

how to add a service:
	git submodule add git@gitlab.com:utopia-planitia/kured.git services/kured
	utopiactl configure kured
	utopiactl deploy kured
	git commit -a -m "added kured (kubernetes reboot daemon)"
	git push origin master

how to update a service:
	utopiactl exec kured git pull
	utopiactl configure kured
	utopiactl deploy kured
	git commit -a -m "updated kured (kubernetes reboot daemon)"
	git push origin master

how to compare current state with available updates:
	utopiactl compare-versions all
`

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to determine current working directory: %v", err)
	}

	if len(os.Args) <= 2 {
		printHelp()
		return
	}

	err = utopia.ExecuteCommandline(cwd, os.Args)
	if err != nil {
		log.Fatalf("command failed: %v", err)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	os.Stdout.WriteString(help)
}
