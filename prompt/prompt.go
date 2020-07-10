package prompt

import (
	"errors"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	errRequiredMissing = errors.New("required missing")
	errInvalidBoolean  = errors.New("invalid boolean")
)

type readerWriter struct {
	io.Reader
	io.Writer
}

func prompt(text string, password bool) (string, error) {
	stdin := os.Stdin
	stdout := os.Stdout
	term := terminal.NewTerminal(readerWriter{stdin, stdout}, text)
	stdinFD := int(stdin.Fd())
	stdinState, err := terminal.MakeRaw(stdinFD)
	if err != nil {
		return "", err
	}
	defer terminal.Restore(stdinFD, stdinState)

	var line string
	if password {
		line, err = term.ReadPassword(text)
	} else {
		line, err = term.ReadLine()
	}
	if err != nil {
		return "", err
	}
	return line, nil
}

func Prompt(text string, required bool) (string, error) {
	line, err := prompt(text, false)
	if err != nil {
		return line, err
	}
	if required && line == "" {
		return line, errRequiredMissing
	}
	return line, err
}

func PromptDefault(text string, dft string) (string, error) {
	line, err := prompt(text, false)
	if err != nil {
		return line, err
	}
	if line == "" {
		return dft, nil
	}
	return line, err
}

func Password(text string) (string, error) {
	return prompt(text, true)
}

func Ask(question string, dft bool) (bool, error) {
	line, err := prompt(question, false)
	if err != nil {
		return false, err
	}
	if line == "" {
		return dft, nil
	}
	return line == "y" || line == "Y" || line == "yes" || line == "Yes" || line == "YES", nil
}
