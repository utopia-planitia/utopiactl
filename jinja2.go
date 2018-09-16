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
{{ range . }}
  - name: render configuration template
    template:
      src: "{{ .Src }}"
      dest: "{{ .Dest }}"
{{ end }}
`
const jinjaSuffix = ".j2"

type jinja2Template struct {
	Src, Dest string
}

func renderJinja2(customizePath string, t []jinja2Template) error {

	playbook, err := ioutil.TempFile(customizePath, playbookPrefix)
	if err != nil {
		return fmt.Errorf("failed to create playbook: %v", err)
	}
	defer os.Remove(playbook.Name())

	playbookTemplate, err := template.New("playbook").Parse(playbookSource)
	if err != nil {
		return fmt.Errorf("failed to create playbook: %v", err)
	}

	err = playbookTemplate.Execute(playbook, t)
	if err != nil {
		return fmt.Errorf("failed to render playbook: %v", err)
	}

	err = playbook.Close()
	if err != nil {
		return fmt.Errorf("failed to close playbook: %v", err)
	}

	c, _ := ioutil.ReadFile(playbook.Name())
	log.Println(playbook.Name())
	log.Printf("%+s", c)

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
