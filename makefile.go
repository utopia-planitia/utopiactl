package utopia

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

func makeConfigure(configTemplatesDir string) error {
	cmd := exec.Command("make", "configure")
	cmd.Dir = configTemplatesDir
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

make deploy: #@setup apply all configurations
	kubectl apply -r -f{{ range .Applys }} {{ . }}{{ end }}{{ range .Makes }}
	$(MAKE) -C {{ . }} deploy{{ end }}

`

func generateMakefile(directory, customizedPath string) error {

	repos, err := subDirectories(customizedPath)
	if err != nil {
		return fmt.Errorf("failed to list repos: %v", err)
	}

	makes := []string{}
	applys := []string{}
	for _, repo := range repos {
		repoPath := filepath.Join(directory, repo, templatesDir)
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
		Makes  []string
		Applys []string
	}{
		Makes:  makes,
		Applys: applys,
	})
	if err != nil {
		return fmt.Errorf("failed to write Makefile: %v", err)
	}

	return nil
}

func subDirectories(path string) ([]string, error) {
	contents, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %v", err)
	}

	subDirectories := []string{}
	for _, content := range contents {
		if !content.IsDir() {
			continue
		}
		subDirectories = append(subDirectories, content.Name())
	}
	return subDirectories, nil
}
