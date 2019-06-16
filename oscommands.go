package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

func CheckIfCommandExists(command string) bool {
	_, err := exec.LookPath(command)
	if err != nil {
		return false
	}
	return true
}

func RunCommand(command []string) {
	c := "sudo"
	cmd := exec.Command(c, command ...)

	outp, err := cmd.StdoutPipe()

	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(outp)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

}