package main

import (
	"fmt"
	i "github.com/giorgiovilardo/crittomane/internal"
	"golang.org/x/term"
	"os"
	"strings"
	"syscall"
)

func main() {
	command, password, err := parseArgs(os.Args)
	if err != nil {
		fmt.Println("Crittomane v1.1.0: pass e or d as the command and optionally a password.")
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
	if len(args) < 2 || len(args) > 3 {
		return "", "", fmt.Errorf("expected up to 2 arguments, got %d", len(args))
	}

	command := strings.ToLower(args[1])

	if command != "e" && command != "d" {
		return "", "", fmt.Errorf("invalid command %q", command)
	}

	var pass string
	if len(args) == 3 {
		pass = args[2]
	} else {
		askedPassword, err := askPassword()
		if err != nil {
			return "", "", fmt.Errorf("error getting password")
		}
		pass = askedPassword
	}

	return command, pass, nil
}

func askPassword() (string, error) {
	fmt.Print("Password: ")
	bytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
