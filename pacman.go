package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// Retrieve all non aur packages
func BasePackages() ([]string, error) {

	pacman := exec.Command("pacman", "-Qent")
	awk := exec.Command("awk", "{print $1}")

	// Retrieve foreign packages and pass as args to grep -v
	var aur []byte
	var err error
	if aur, err = exec.Command("pacman", "-Qm").Output(); err != nil {
		return nil, err
	}
	args := fmt.Sprintf(" '%s'", string(aur))
	vgrep := exec.Command("grep", "-v", args)

	if awk.Stdin, err = pacman.StdoutPipe(); err != nil {
		return nil, err
	}
	if vgrep.Stdin, err = awk.StdoutPipe(); err != nil {
		return nil, err
	}

	if err = pacman.Start(); err != nil {
		return nil, err
	}
	if err = awk.Start(); err != nil {
		return nil, err
	}

	output, err := vgrep.Output()
	if err != nil {
		return nil, err
	}

	text := strings.TrimSpace(string(output))

	res := strings.Split(text, "\n")
	return res, nil
}
