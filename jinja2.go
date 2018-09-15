package utopia

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alecthomas/template"
)

const playbookPrefix = "playbook"
const playbookSource = `
---
- hosts: localhost
  tasks:
  - name: render configuration template
    template:
      src: "{{ .Src }}"
      dest: "{{ .Dest }}"
`
const jinjaSuffix = ".j2"

func renderJinja2(customizePath, src, dest string) error {

	src, err := filepath.Abs(src)
	if err != nil {
		return fmt.Errorf("failed to get absolute of src: %v", err)
	}

	playbook, err := ioutil.TempFile(customizePath, playbookPrefix)
	if err != nil {
		return fmt.Errorf("failed to create playbook: %v", err)
	}
	defer os.Remove(playbook.Name())

	playbookTemplate, err := template.New("playbook").Parse(playbookSource)
	if err != nil {
		return fmt.Errorf("failed to create playbook: %v", err)
	}

	err = playbookTemplate.Execute(playbook, struct {
		Src, Dest string
	}{
		Src:  src,
		Dest: strings.TrimSuffix(dest, jinjaSuffix),
	})
	if err != nil {
		return fmt.Errorf("failed to render playbook: %v", err)
	}

	err = playbook.Close()
	if err != nil {
		return fmt.Errorf("failed to close playbook: %v", err)
	}

	return executeAnsiblePlaybook(playbook.Name())
}

func executeAnsiblePlaybook(playbook string) error {
	cmd := exec.Command("ansible-playbook", playbook)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("ansible output: %s", output)
		return fmt.Errorf("failed to execute ansible playbook: %v", err)
	}
	return nil
}
