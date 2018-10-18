package utopia

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Deploy applies the application and configuration definition to kubernetes.
func Deploy(directory string, services []string) error {

	if len(services) == 0 {
		err := makeDeploy(directory)
		if err != nil {
			return fmt.Errorf("clsuter deployment failed: %v", err)
		}
		return nil
	}

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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("make deploy: %v", err)
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("kubectl apply: %v", err)
	}
	return nil
}
