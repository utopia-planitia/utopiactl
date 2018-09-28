package utopia

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

func makeConfigure(generatedConfigDir string) error {
	cmd := exec.Command("make", "configure")
	cmd.Dir = generatedConfigDir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PWD=%s", generatedConfigDir),
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("make output: %s", output)
		return fmt.Errorf("failed to execute make configure: %v", err)
	}
	return nil
}

const makefileSource = `
include ../kubernetes/etc/help.mk
include ../kubernetes/etc/cli.mk

deploy: applications configurations ##@setup apply all applications and configurations

applications: ##@setup apply all applications{{ range .Applications }}
	cd ../{{ . }} && make deploy{{ end }}

configurations: ##@setup apply all configurations
	$(CLI) kubectl apply -R \
{{ range .Applys }}		-f {{ . }} \
{{ end }}
{{ range .Makes }}	cd {{ . }} && make deploy
{{ end }}`

func generateMakefile(directory, customizedPath string) error {

	repos, err := subDirectories(directory)
	if err != nil {
		return fmt.Errorf("failed to list repos: %v", err)
	}

	makes := []string{}
	applys := []string{}
	applications := []string{}
	for _, repo := range repos {

		if repo == customizedRepo || repo == "hetzner" || repo == "kubernetes" {
			continue
		}

		if _, err := os.Stat(filepath.Join(directory, repo, "Makefile")); err == nil {
			applications = append(applications, repo)
		}

		repoPath := filepath.Join(directory, repo, templatesDir)
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			continue
		}

		if _, err := os.Stat(filepath.Join(repoPath, "Makefile")); err == nil {
			makes = append(makes, repo)
		} else {
			applys = append(applys, repo)
		}
	}

	makefileTemplate, err := template.New("makefile").Parse(makefileSource)
	if err != nil {
		return fmt.Errorf("failed to parse makefile template: %v", err)
	}

	makefile, err := os.OpenFile(filepath.Join(customizedPath, "Makefile"), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("opening Makefile failed: %v", err)
	}
	defer makefile.Close()
	err = makefile.Truncate(0)
	if err != nil {
		return fmt.Errorf("resetting Makefile failed: %v", err)
	}

	err = makefileTemplate.Execute(makefile, struct {
		Makes        []string
		Applys       []string
		Applications []string
	}{
		Makes:        makes,
		Applys:       applys,
		Applications: applications,
	})
	if err != nil {
		return fmt.Errorf("failed to write Makefile: %v", err)
	}

	return nil
}
