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
all: {{ if .HetznerExists }}hetzner {{ end }}{{- if .KubernetesExists -}}kubernetes {{ end }}services configurations ##@setup deploy everything

{{ if .KubernetesExists -}}
.PHONY: kubernetes
kubernetes: ##@setup deploy kubernetes
	cd services/kubernetes && make deploy

{{ end -}}{{ if .HetznerExists -}}
.PHONY: hetzner
hetzner: ##@setup run maintenance for hetzner nodes
	cd services/hetzner && make maintenance

{{ end -}}
.PHONY: deploy
deploy: services configurations ##@setup apply all applications and configurations

.PHONY: services
services: ##@setup apply all applications{{ range .Applications }}
	cd services/{{ . }} && make deploy{{ end }}

.PHONY: configurations
configurations: ##@setup apply all configurations
{{- if .Applys }}
	$(CLI) kubectl apply -R \
{{ range .Applys }}		-f configurations/{{ . }} \
{{ end }}
{{- end }}
{{ range .Makes }}	cd configurations/{{ . }} && make deploy
{{ end }}`

func generateMakefile(directory string) error {

	services := subDirectories(filepath.Join(directory, "services"))
	services = moveServiceToFirst(services, "metrics")
	services = moveServiceToFirst(services, "storage")
	services = moveServiceToFirst(services, "priority-class-patcher")

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

	configs := subDirectories(filepath.Join(directory, "configurations"))
	configs = moveServiceToFirst(configs, "metrics")
	configs = moveServiceToFirst(configs, "storage")

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

	hetznerExists := false
	if _, err := os.Stat(filepath.Join(directory, "services", "hetzner", "Makefile")); err == nil {
		hetznerExists = true
	}

	kubernetesExists := false
	if _, err := os.Stat(filepath.Join(directory, "services", "kubernetes", "Makefile")); err == nil {
		kubernetesExists = true
	}

	err = makefileTemplate.Execute(makefile, struct {
		Makes            []string
		Applys           []string
		Applications     []string
		HetznerExists    bool
		KubernetesExists bool
	}{
		Makes:            cfgMakes,
		Applys:           cfgApplys,
		Applications:     applications,
		HetznerExists:    hetznerExists,
		KubernetesExists: kubernetesExists,
	})
	if err != nil {
		return fmt.Errorf("failed to write Makefile: %v", err)
	}

	return nil
}

func moveServiceToFirst(services []string, service string) []string {
	if !contains(services, service) {
		return services
	}

	svcs := []string{service}
	for _, s := range services {
		if s == service {
			continue
		}
		svcs = append(svcs, s)
	}
	return svcs
}
