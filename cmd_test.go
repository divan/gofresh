package main

import "testing"

func TestCmd(t *testing.T) {
	data := []byte(`line1
line2
line3`)

	lines := bytes2strings(data)
	if len(lines) != 3 {
		t.Fatal("Expecting 3 lines, but have", len(lines))
	}
}
