package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-cmd/cmd"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
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

// ExecuteStream : run a bash command and stream output
func ExecuteStream(command string, cwd string, fn func(stdout string, stderr string)) (commandExe *cmd.Cmd) {
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}

	exeCmd := cmd.NewCmdOptions(cmdOptions, "bash", "-c", command)

	go func() {
		for {
			select {
			case line := <-exeCmd.Stdout:
				fn(line, "")
			case line := <-exeCmd.Stderr:
				fn("", line)
			}
		}
	}()

	<-exeCmd.Start()

	for len(exeCmd.Stdout) > 0 || len(exeCmd.Stderr) > 0 {
		time.Sleep(10 * time.Millisecond)
	}

	return exeCmd
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

// ExecuteCwdStreamWithEnv : stream a bash command
func ExecuteCwdStreamWithEnv(command string, cwd string, envVars map[string]string, fn func(stdout string)) (outStr string, errStr string, exitCode error) {
	cwd = ResolvePath(cwd)
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = cwd

	cmd.Env = os.Environ()
	for envKey, envVal := range envVars {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", envKey, envVal))
	}

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
