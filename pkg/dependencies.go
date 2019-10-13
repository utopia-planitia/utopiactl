package utopia

import (
	"bytes"
	"os"
	"path/filepath"
)

// Dependencies generates a graph combining the dependencies between services.
func Dependencies(directory string, services []string) error {

	var buffer bytes.Buffer

	buffer.WriteString("digraph {\n")
	buffer.WriteString("define(digraph,subgraph)\n")

	for _, svc := range services {
		file := filepath.Join("services", svc, "dependencies.dot")
		if !fileExists(filepath.Join(directory, file)) {
			continue
		}
		buffer.WriteString("include(" + file + ")\n")
	}
	buffer.WriteString("}\n")

	os.Stdout.WriteString(buffer.String())

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
