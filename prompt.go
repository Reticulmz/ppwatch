package main

import (
	"fmt"
	"gopkg.in/readline.v1"
)

func PromptForInput(prompt, defaultvalue string) (string, error) {
	var promptstr string
	if defaultvalue != "" {
		promptstr = fmt.Sprintf("(press enter for default: %s) > ", defaultvalue)
	} else {
		promptstr = "> "
	}

	rl, err := readline.New(promptstr)
	if err != nil {
		return "", err
	}

	defer rl.Close()

	fmt.Printf("\n%s\n\n", prompt)
	line, err := rl.Readline()
	if err != nil {
		// EOF or interrupt: return default
		return defaultvalue, nil
	}

	if line == "" {
		return defaultvalue, nil
	}

	return line, nil
}
