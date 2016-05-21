package main

import (
	"testing"
	//"os"
)

func TestMainRun(t testing.T) {
	cases := []struct {
		name string
		cmd  string
		arg1 string
		arg2 string
	}{
		{"cfs", "run", "ls", "-l"},
		{"cfs", "child", "ls", "-la"},
		{"cfs", "invalid", "no", "-n"},
		{"cfs", "run", "invalidArg", "invalidFlag"},
	}

	app := makeCmd()

	for _, entry := range cases {
		app.Run(entry.name, entry.cmd, entry.arg1, entry.arg2)

	}

}
