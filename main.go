package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		wd, err := os.Getwd()
		sendToStdErr(err)

		dirs := strings.Split(wd, "/")
		curDir := dirs[len(dirs)-1]

		host, err := os.Hostname()
		sendToStdErr(err)

		user, err := user.Current()
		sendToStdErr(err)

		fmt.Printf("[%s@%s %s]> ", user.Username, host, curDir)

		input, err := reader.ReadString('\n')
		sendToStdErr(err)

		if err = execInput(input); err != nil {
			sendToStdErr(err)
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

		fmt.Println(dirToCD)
		return os.Chdir(dirToCD)
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	lastCDedDir = wd

	return os.Chdir(args[1])
}

func sendToStdErr(err error) {
	fmt.Fprintln(os.Stderr, err)
}
