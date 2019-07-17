package utopia

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func VerifyTests(directory string, services []string, testAllServices bool) ([]error, error) {

	if len(services) == 0 {
		return []error{}, fmt.Errorf("service list is missing")
	}

	failedServices := []error{}

	for _, svc := range services {

		hasTests, err := makeTargetExist("tests", directory, svc)
		if err != nil {
			return []error{}, fmt.Errorf("failed to search 'make tests': %v", err)
		}

		if !hasTests {
			log.Printf("service %s does not have tests\n", svc)
			continue
		}

		log.Printf("execute command for service %s\n", svc)
		err = execCommand(filepath.Join(directory, "services", svc), []string{"make", "-s", "tests"})
		if err != nil {
			if !testAllServices {
				return []error{}, fmt.Errorf("service %s tests failed: %v", svc, err)
			}
			failedServices = append(failedServices, fmt.Errorf("service %s tests failed: %v", svc, err))
		}
	}

	if len(failedServices) != 0 {
		return failedServices, fmt.Errorf("some services had errors")
	}
	return []error{}, nil
}

func makeTargetExist(target, directory, service string) (bool, error) {
	err := execCommandSilent(filepath.Join(directory, "services", service), []string{"make", "-q", target})
	code, err := exitCode(err)
	if err != nil {
		return false, fmt.Errorf("failed to search make target: %v", err)
	}
	return code == 1, nil
}

func execCommandSilent(dir string, command []string) error {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PWD=%s", dir),
		"DOCKER_INTERACTIVE= ",
		"MAKE=make",
	)
	cmd.Dir = dir
	return cmd.Run()
}

func exitCode(err error) (int, error) {
	if err == nil {
		return 0, nil
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode(), nil
	}
	return 0, fmt.Errorf("failed to extract exitcode")
}
