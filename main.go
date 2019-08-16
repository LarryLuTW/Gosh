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
	args := strings.Split(input, " ")
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

	fmt.Printf("%s@%s %s > ", u.Username, host, wd)
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
