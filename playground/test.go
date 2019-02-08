package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func main() {
	for i := 1; i < 3; i++ {
		go runShellCommand("./test.sh")
	}

	time.Sleep(60000 * time.Millisecond)
}

func runShellCommand(commandText string) {
	cmd := exec.Command(commandText)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
