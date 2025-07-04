package main

import (
	"bufio"
	"errors"
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
		handleErr(err)

		dirs := strings.Split(wd, "/")
		curDir := dirs[len(dirs)-1]

		host, err := os.Hostname()
		handleErr(err)

		user, err := user.Current()
		handleErr(err)

		fmt.Printf("[%s@%s %s]> ", user.Username, host, curDir)

		input, err := reader.ReadString('\n')
		handleErr(err)

		if err = execInput(input); err != nil {
			handleErr(err)
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
		updateLastCDedDir()

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		return os.Chdir(homeDir)
	}

	if args[1] == "-" {
		if len(lastCDedDir) == 0 {
			return errors.New("previous working directory not set")
		}

		dirToCD := lastCDedDir

		updateLastCDedDir()

		fmt.Println(dirToCD)
		return os.Chdir(dirToCD)
	}

	updateLastCDedDir()

	return os.Chdir(args[1])
}

func updateLastCDedDir() {
	wd, err := os.Getwd()
	if err != nil {
		handleErr(err)
	}

	lastCDedDir = wd
}

func handleErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
