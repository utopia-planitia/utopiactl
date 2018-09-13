package utopia

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/alecthomas/template"
)

const playbookPrefix = "playbook"
const playbookSource = `
---
- hosts: localhost
  tasks:
	- name: render configuration template
	  template:
		src: "{{ src }}"
		dest: "{{ dest }}"
`

func renderJinja2(customizePath, path, dest string) error {

	playbook, err := ioutil.TempFile(os.TempDir(), playbookPrefix)
	if err != nil {
		return fmt.Errorf("failed to create playbook: %v", err)
	}
	defer os.Remove(playbook.Name())

	playbookTemplate, err := template.New("playbook").Parse(playbookSource)
	if err != nil {
		return fmt.Errorf("failed to create playbook: %v", err)
	}

	err = playbookTemplate.Execute(playbook, struct {
		src, dest string
	}{
		src:  path,
		dest: dest,
	})
	if err != nil {
		return fmt.Errorf("failed to render playbook: %v", err)
	}

	err = playbook.Close()
	if err != nil {
		return fmt.Errorf("failed to close playbook: %v", err)
	}

	return executeAnsiblePlaybook(customizePath, playbook.Name())
}

func executeAnsiblePlaybook(customizePath, playbook string) error {
	cmd := exec.Command("ansible-playbook")
	cmd.Dir = customizePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(output)
		return fmt.Errorf("failed to execute ansible playbook: %v", err)
	}
	return nil
}
