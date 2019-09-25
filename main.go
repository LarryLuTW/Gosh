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

func parseArgs(input string) []string {
	if strings.HasPrefix(input, "alias") {
		return strings.SplitN(input, " ", 2)
	}
	return strings.Split(input, " ")
}

func executeInput(input string) error {
	input = os.ExpandEnv(input)
	input = expandWildcardInCmd(input)

	if strings.HasPrefix(input, `\`) {
		input = input[1:]
	} else {
		input = expandAlias(input)
	}

	args := parseArgs(input)

	if args[0] == "cd" {
		err := os.Chdir(args[1])
		return err
	}

	if args[0] == "alias" {
		// args = ["alias", "ls='ls -l'"]
		// kv = ["ls", "'ls -l'"]
		kv := strings.Split(args[1], "=")
		// key = "ls", val = "ls -l"
		key, val := kv[0], strings.Trim(kv[1], "'")
		setAlias(key, val)
		return nil
	}

	if args[0] == "unalias" {
		// args = ["unalias", "ls"]
		key := args[1]
		unsetAlias(key)
		return nil
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
