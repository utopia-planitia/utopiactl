package utopia

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const makefileSource = `
include services/kubernetes/etc/help.mk
include services/kubernetes/etc/cli.mk

.PHONY: all
all: hetzner kubernetes services configurations ##@setup deploy everything

.PHONY: kubernetes
kubernetes: ##@setup deploy kubernetes
	cd services/kubernetes && make deploy

.PHONY: hetzner
hetzner: ##@setup run maintenance for hetzner nodes
	cd services/hetzner && make maintenance

.PHONY: deploy
deploy: services configurations ##@setup apply all applications and configurations

.PHONY: services
services: ##@setup apply all applications{{ range .Applications }}
	cd services/{{ . }} && make deploy{{ end }}

.PHONY: configurations
configurations: ##@setup apply all configurations
	$(CLI) kubectl apply -R \
{{ range .Applys }}		-f configurations/{{ . }} \
{{ end }}
{{ range .Makes }}	cd configurations/{{ . }} && make deploy
{{ end }}`

func generateMakefile(directory string) error {

	services, err := subDirectories(filepath.Join(directory, "services"))
	if err != nil {
		return fmt.Errorf("failed to list services: %v", err)
	}

	applications := []string{}
	for _, svc := range services {
		if contains([]string{"hetzner", "kubernetes"}, svc) {
			continue
		}

		if _, err := os.Stat(filepath.Join(directory, "services", svc, "Makefile")); err != nil {
			continue
		}
		applications = append(applications, svc)
	}

	configs, err := subDirectories(filepath.Join(directory, "configurations"))
	if err != nil {
		return fmt.Errorf("failed to list services: %v", err)
	}

	cfgMakes := []string{}
	cfgApplys := []string{}
	for _, cfg := range configs {
		if contains([]string{"hetzner", "kubernetes"}, cfg) {
			continue
		}

		if _, err := os.Stat(filepath.Join(directory, "configurations", cfg, "Makefile")); err == nil {
			cfgMakes = append(cfgMakes, cfg)
			continue
		}
		cfgApplys = append(cfgApplys, cfg)
	}

	makefileTemplate, err := template.New("makefile").Parse(makefileSource)
	if err != nil {
		return fmt.Errorf("failed to parse makefile template: %v", err)
	}

	makefile, err := os.OpenFile(filepath.Join(directory, "Makefile"), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("opening Makefile failed: %v", err)
	}
	defer makefile.Close()
	err = makefile.Truncate(0)
	if err != nil {
		return fmt.Errorf("truncating Makefile failed: %v", err)
	}

	err = makefileTemplate.Execute(makefile, struct {
		Makes        []string
		Applys       []string
		Applications []string
	}{
		Makes:        cfgMakes,
		Applys:       cfgApplys,
		Applications: applications,
	})
	if err != nil {
		return fmt.Errorf("failed to write Makefile: %v", err)
	}

	return nil
}
