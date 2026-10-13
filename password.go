package main

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
)

// Not a flag: a password given as an argument stays in the shell history
// and in the process list.
func readPassword() (string, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return "", errors.New("the password is read from a terminal, and stdin is not one")
	}
	fmt.Fprint(os.Stderr, "Password: ")
	password, err := term.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)
	return string(password), err
}
