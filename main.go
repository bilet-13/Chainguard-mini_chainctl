package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"mychainctl/cmd"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		lower := strings.ToLower(line)
		if lower == "exit" || lower == "quit" {
			break
		}

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		if args[0] == "mychainctl" {
			args = args[1:]
		}
		if len(args) == 0 {
			continue
		}

		if err := cmd.ExecuteWithArgs(args); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
