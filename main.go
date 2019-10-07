package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strings"
	"syscall"
)

func parseArgs(input string) []string {
	if strings.HasPrefix(input, "alias") {
		return strings.SplitN(input, " ", 2)
	}
	return strings.Split(input, " ")
}

func lookCommand(cmd string) {
	value := aliasTable[cmd]
	if value != "" {
		fmt.Printf("%s: aliased to %s\n", cmd, value)
		return
	}

	value, err := exec.LookPath(cmd)
	if err == nil {
		fmt.Printf("%s: %s\n", cmd, value)
		return
	}

	fmt.Printf("%s NOT FOUND\n", cmd)
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

	if args[0] == "which" {
		for _, cmd := range args[1:] {
			lookCommand(cmd)
		}
		return nil
	}

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

	inputStream := os.Stdin

	if len(args) > 2 && args[len(args)-2] == "<" {
		filename := args[len(args)-1]
		file, err := os.Open(filename)
		if err != nil {
			return err
		}

		inputStream = file
		args = args[:len(args)-2]
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdin = inputStream
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

func initialize() {
	homeDir, _ := os.UserHomeDir()
	file, err := os.Open(homeDir + "/.goshrc")
	if err != nil {
		return
	}

	goshrcReader := bufio.NewReader(file)

	for {
		input, err := goshrcReader.ReadString('\n')
		if err == io.EOF {
			return
		}

		input = strings.TrimSpace(input)
		executeInput(input)
	}
}

func main() {
	initialize()

	stdin := bufio.NewReader(os.Stdin)

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGINT)

	handleSignals := func() {
		for {
			sig := <-signalCh
			fmt.Println("Received signal:", sig)
		}
	}

	go handleSignals()

	for {
		showPrompt()

		input, _ := stdin.ReadString('\n')
		input = strings.TrimSpace(input)

		err := executeInput(input)
		if err != nil {
			fmt.Println(err)
		}
	}
}
