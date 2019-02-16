package utils

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

// Execute : run a bash command
func Execute(command string) string {
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = "/"
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return out.String()
}

// ExecuteCwd : run a bash command
func ExecuteCwd(command string, cwd string) string {
	cwd = ResolvePath(cwd)
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = cwd
	var out bytes.Buffer
	var outerr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outerr
	err := cmd.Run()
	if err != nil {
		fmt.Printf(out.String())
		fmt.Printf(outerr.String())
		log.Fatal(err)
	}
	return out.String()
}
