package utils

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// ResolvePath : Resolve a path to include users home directory
func ResolvePath(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	if path == "~" {
		// In case of "~", which won't be caught by the "else if"
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		path = filepath.Join(dir, path[2:])
	}
	return path
}

// EnsureDirectoryExists : Create a directory if it doesn't already exist
func EnsureDirectoryExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}
