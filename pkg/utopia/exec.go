package utopia

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Exec executes a given command in every service.
func Exec(directory string, services []string, command []string) error {

	for _, svc := range services {

		log.Printf("execute command for service %s", svc)

		cmd := exec.Command(command[0], command[1:]...)
		cmd.Dir = filepath.Join(directory, "services", svc)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("command throw error: %v", err)
		}
		_, err = os.Stdout.Write(output)
		if err != nil {
			log.Printf("failed to print output: %v", err)
		}
	}

	return nil
}
