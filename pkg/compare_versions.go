package utopia

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func CompareVersions(directory string, services []string) error {

	if len(services) == 0 {
		return fmt.Errorf("service list is missing")
	}

	for _, svc := range services {

		hash, err := currentHash(directory, svc)
		if err != nil {
			return fmt.Errorf("failed to get hash of submodule: %v", err)
		}
		c, err := CompareVersionsDelta(directory, svc)
		if err != nil {
			return fmt.Errorf("failed to count changes: %v", err)
		}
		fmt.Printf("service %s (%d)\n", svc, c)

		cmd := []string{"git", "log", "--graph", "--oneline", "--decorate=false", fmt.Sprintf("%s..origin/master", hash)}
		err = execCommand(filepath.Join(directory, "services", svc), cmd)
		if err != nil {
			return fmt.Errorf("git log failed: %v", err)
		}
	}

	return nil
}

var re = regexp.MustCompile(`[0-9a-fA-F]{40}`)

func currentHash(dir string, svc string) (string, error) {
	cmd := exec.Command("git", "ls-tree", "HEAD", "services/"+svc)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PWD=%s", dir),
	)
	cmd.Dir = dir
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command throw error: %v", err)
	}

	hash := re.Find(stdoutStderr)
	if hash == nil {
		return "", fmt.Errorf("could not detect submodule hash for service %s", svc)
	}

	return string(hash), nil
}

func CompareVersionsDelta(directory string, svc string) (int, error) {
	hash, err := currentHash(directory, svc)
	if err != nil {
		return 0, fmt.Errorf("failed to get hash of submodule: %v", err)
	}

	command := []string{"git", "log", "--oneline", "--decorate=false", fmt.Sprintf("%s..origin/master", hash)}

	dir := filepath.Join(directory, "services", svc)
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PWD=%s", dir),
	)
	cmd.Dir = dir
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("command throw error: %v", err)
	}

	if len(stdoutStderr) == 0 {
		return 0, nil
	}

	return bytes.Count(stdoutStderr, []byte("\n")), nil
}
