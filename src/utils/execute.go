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
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = cwd
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf(out.String())
		log.Fatal(err)
	}
	return out.String()
}
