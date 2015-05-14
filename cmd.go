package main

import (
	"os/exec"
	"strings"
)

func Run(command, dir string) ([]string, error) {
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
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
