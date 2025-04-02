package main

import (
	"fmt"
	i "github.com/giorgiovilardo/crittomane/internal"
	"os"
	"strings"
)

func main() {
	command, password, err := parseArgs(os.Args)
	if err != nil {
		fmt.Println("Wrong argument: pass e or d as the command, then a password.")
		os.Exit(1)
	}

	switch command {
	case "e":
		zipBuffer, err := i.TarPrivate()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		encrypted, err := i.EncryptBytes(zipBuffer.Bytes(), password)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = os.WriteFile("private.ctm", encrypted, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Done!")
		os.Exit(0)
	case "d":
		file, err := os.ReadFile("private.ctm")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		decrypted, err := i.DecryptBytes(file, password)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = i.UntarBytes(decrypted)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Done!")
		os.Exit(0)
	}
}

func parseArgs(args []string) (string, string, error) {
	if len(args) != 3 {
		return "", "", fmt.Errorf("expected 2 arguments, got %d", len(args))
	}

	command := strings.ToLower(args[1])
	pass := strings.ToLower(args[2])

	if command != "e" && command != "d" {
		return "", "", fmt.Errorf("invalid command %q", command)
	}

	return command, pass, nil
}
