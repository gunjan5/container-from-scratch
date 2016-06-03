package main

import (
	"testing"
	//"os"
)

func TestMainRun(t *testing.T) {
	cases := []struct {
		name  string
		state string
		image string
		cmd   string
	}{
		{"cfs", "run", "TinyCore", "ls"},
		{"cfs", "run", "SlitazOS", "whoami"},
		{"cfs", "run", "BusyBox", "pwd"},
		//	{"cfs", "run", "pwd", ""},
		//	{"cfs", "invalid", "no", "-n"},
		//	{"cfs", "run", "invalidArg", "invalidFlag"},
	}

	app := makeCmd()

	//TODO: figure out how to test the output for this :/

	for _, entry := range cases {
		app.Run([]string{entry.name, entry.state, entry.image, entry.cmd})

	}

}
