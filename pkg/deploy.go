package utopia

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Deploy applies the application and configuration definition to kubernetes.
func Deploy(directory string, services []string) error {

	for _, svc := range services {
		err := makeDeploy(filepath.Join(directory, "services", svc))
		if err != nil {
			return fmt.Errorf("application deployment failed for service %v: %v", svc, err)
		}
		err = deployConfiguration(directory, svc)
		if err != nil {
			return fmt.Errorf("application deployment failed for service %v: %v", svc, err)
		}
	}
	return nil
}

func makeDeploy(dir string) error {
	cmd := exec.Command("make", "deploy")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PWD=%s", dir),
		"DOCKER_INTERACTIVE= ",
	)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("make deploy: %v: %s", err, output)
	}
	return nil
}

func deployConfiguration(directory string, svc string) error {
	if _, err := os.Stat(filepath.Join(directory, "configurations", svc, "Makefile")); err == nil {
		return makeDeploy(filepath.Join(directory, "configurations", svc))
	}
	return applyConfiguration(directory, svc)
}

func applyConfiguration(directory string, svc string) error {
	cmd := exec.Command("make", "cli")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PWD=%s", directory),
		"DOCKER_INTERACTIVE= ",
		fmt.Sprintf("CMD=kubectl apply -R -f %s", filepath.Join("configurations", svc)),
	)
	cmd.Dir = directory
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("kubectl apply: %v: %v", err, output)
	}
	return nil
}
