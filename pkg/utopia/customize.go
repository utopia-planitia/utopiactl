package utopia

import (
	"fmt"
)

const templatesDir = "config-templates"

// Customize updates the customized cluster repository. All services have to be
// located in the services directory, only repositories listed in repos are
// updated. It renders jinja2 templates using Ansible and creates a Makefile to
// apply the custom configuration to Kubernetes.
// A Makefile in the config-templates directory can extend the default template
// & kubectl behavior. Makefile targets 'make configure' and 'make deploy' are
// called.
func Customize(directory string, repos []string) error {

	err := copyAnsibleVars(directory, repos)
	if err != nil {
		return fmt.Errorf("collection ansible variables failed: %v", err)
	}

	err = autogenerateConfigs(directory, repos)
	if err != nil {
		return fmt.Errorf("config autogeneration failed: %v", err)
	}

	err = generateMakefile(directory)
	if err != nil {
		return fmt.Errorf("failed to generate Makefile: %v", err)
	}

	return nil
}
