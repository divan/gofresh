package main

import (
	"os/exec"
	"strings"
)

func Run(dir, command string, args ...string) ([]string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return bytes2strings(out), nil
}

func bytes2strings(data []byte) []string {
	isNewline := func(r rune) bool {
		return r == '\n'
	}
	return strings.FieldsFunc(string(data), isNewline)
}
