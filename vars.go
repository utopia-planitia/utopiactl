package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
)

func vars(repoPath, customizePath string) error {
	err := copyIfExists(repoPath, customizePath, "host_vars")
	if err != nil {
		return fmt.Errorf("failed to copy host_vars: %v", err)
	}
	copyIfExists(repoPath, customizePath, "group_vars")
	if err != nil {
		return fmt.Errorf("failed to copy group_vars: %v", err)
	}
	return nil
}

func copyIfExists(source, target, path string) error {
	src := filepath.Join(source, path)
	log.Printf("src %v \n", src)

	stat, err := os.Stat(src)
	if os.IsNotExist(err) {
		return nil
	}
	if !stat.IsDir() {
		return fmt.Errorf("%v is not a directory", err)
	}

	dest := filepath.Join(target, path)
	log.Printf("%v %v \n", src, dest)

	return copy.Copy(src, dest)
}
