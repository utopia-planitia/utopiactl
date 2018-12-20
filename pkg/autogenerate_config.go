package utopia

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func autogenerateConfigs(directory string, services []string) error {
	jt := []jinja2Template{}

	directory, err := filepath.Abs(directory)
	if err != nil {
		return fmt.Errorf("failed to get absolute of directory: %v", err)
	}

	for _, svc := range services {
		src := filepath.Join(directory, "services", svc, "config-templates")
		dest := filepath.Join(directory, "configurations", svc)
		if _, err := os.Stat(src); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(src, walkConfig(&jt, src, dest))
		if err != nil {
			return fmt.Errorf("customization failed for service %v: %v", svc, err)
		}
	}

	ansiblePath := filepath.Join(directory, "ansible")
	err = renderJinja2(ansiblePath, jt)
	if err != nil {
		return fmt.Errorf("jinja2 rendering via ansible failed: %v", err)
	}

	for _, svc := range services {
		dest := filepath.Join(directory, "configurations", svc)
		if _, err := os.Stat(filepath.Join(dest, "Makefile")); err == nil {
			err = makeConfigure(dest)
			if err != nil {
				return fmt.Errorf("make configure (%v): %v", svc, err)
			}
		}
	}

	return nil
}

func walkConfig(jinja2Templates *[]jinja2Template, serviceDir, configDir string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dest := filepath.Join(configDir, strings.TrimPrefix(path, serviceDir))

		if info.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		if filepath.Ext(info.Name()) != jinjaSuffix {
			return copy(path, dest)
		}

		if strings.HasPrefix(strings.TrimPrefix(path, serviceDir), "/roles") {
			return copy(path, dest)
		}

		*jinja2Templates = append(*jinja2Templates, jinja2Template{
			Src:  path,
			Dest: strings.TrimSuffix(dest, jinjaSuffix),
		})

		return nil
	}
}

func copy(src, dest string) error {

	info, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("could not stat %s: %s", src, err)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create file %s: %s", dest, err)
	}
	defer f.Close()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return fmt.Errorf("could not set permissions for file %s: %s", f.Name(), err)
	}

	s, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file %s: %s", src, err)
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	if err != nil {
		return fmt.Errorf("failed to copy from source to destination: %s", err)
	}

	return nil
}

func makeConfigure(generatedConfigDir string) error {
	cmd := exec.Command("make", "configure")
	cmd.Dir = generatedConfigDir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PWD=%s", generatedConfigDir),
		"DOCKER_INTERACTIVE= ",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("make output: %s", output)
		return fmt.Errorf("failed to execute make configure: %v", err)
	}
	return nil
}
