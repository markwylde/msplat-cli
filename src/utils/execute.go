package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
)

// HandleIoError : Handle errors for functions in this file
func HandleIoError(stdout string, stderr string, exit error) {
	if exit != nil {
		fmt.Printf(stdout)
		fmt.Printf(stderr)
		log.Fatalf("exec ended with code %s", exit)
	}
}

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

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
func ExecuteCwd(command string, cwd string) (stdout string, stderr string, exit error) {
	cwd = ResolvePath(cwd)
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = cwd
	var out bytes.Buffer
	var outerr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outerr
	err := cmd.Run()
	return out.String(), outerr.String(), err
}

// ExecuteCwdStream : stream a bash command
func ExecuteCwdStream(command string, cwd string, fn func(stdout string)) (outStr string, errStr string, exitCode error) {
	cwd = ResolvePath(cwd)
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = cwd

	stdout, _ := cmd.StdoutPipe()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fn(m)
	}

	err := cmd.Wait()

	return "", stderr.String(), err
}
