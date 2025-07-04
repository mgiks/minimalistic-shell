package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if err = execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func execInput(input string) error {
	input = strings.TrimSuffix(input, "\n")

	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		return handleCD(args)
	case "exit":
		os.Exit(0)
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

var lastCDedDir string

func handleCD(args []string) error {
	if len(args) < 2 {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		lastCDedDir = wd
		fmt.Println(lastCDedDir)

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		return os.Chdir(homeDir)
	}

	if args[1] == "-" {
		dirToCD := lastCDedDir

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		lastCDedDir = wd

		return os.Chdir(dirToCD)
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	lastCDedDir = wd

	return os.Chdir(args[1])
}
