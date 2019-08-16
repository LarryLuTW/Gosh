package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		input, _ := stdin.ReadString('\n')
		input = strings.TrimSpace(input)

		fmt.Println(input)
	}
}
