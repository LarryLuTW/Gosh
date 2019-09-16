package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func executeInput(input string) error {
	input = os.ExpandEnv(input)

	args := strings.Split(input, " ")

	if args[0] == "cd" {
		err := os.Chdir(args[1])
		return err
	}

	if args[0] == "export" {
		// args = ["export", "FOO=bar"]
		kv := strings.Split(args[1], "=")
		key, val := kv[0], kv[1]
		err := os.Setenv(key, val)
		return err
	}

	if args[0] == "unset" {
		// args = ["unset", "FOO"]
		err := os.Unsetenv(args[1])
		return err
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}

func showPrompt() {
	u, _ := user.Current()
	host, _ := os.Hostname()
	wd, _ := os.Getwd()

	userAndHost := blue(u.Username + "@" + host)
	wd = yellowWithBlueBG(wd)

	fmt.Printf("%s %s > ", userAndHost, wd)
}

func main() {
	stdin := bufio.NewReader(os.Stdin)

	for {
		showPrompt()

		input, _ := stdin.ReadString('\n')
		input = strings.TrimSpace(input)

		err := executeInput(input)
		if err != nil {
			log.Println(err)
		}
	}
}
