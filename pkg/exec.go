package utopia

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Exec executes a given command in every service.
func Exec(directory string, services []string, command []string) error {

	if len(services) == 0 {
		log.Printf("execute command for cluster\n")
		err := execCommand(directory, command)
		if err != nil {
			return fmt.Errorf("cluster command failed: %v", err)
		}
		return nil
	}

	for _, svc := range services {
		log.Printf("execute command for service %s\n", svc)
		err := execCommand(filepath.Join(directory, "services", svc), command)
		if err != nil {
			return fmt.Errorf("service command failed: %v", err)
		}
	}

	return nil
}

func execCommand(dir string, command []string) error {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PWD=%s", dir),
		"DOCKER_INTERACTIVE= ",
	)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command throw error: %v", err)
	}
	return nil
}
