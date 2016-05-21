package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/gunjan5/container-from-scratch/cmd"
)

func main() {

	app := makeCmd()
	app.Run(os.Args)

}

func makeCmd() *cli.App {

	app := cli.NewApp()
	app.Name = "greet"
	app.Usage = "meowwwwwwww out loud"

	fmt.Println(os.Args)

	app.Commands = []cli.Command{
		{
			Name:        "run",
			ShortName:   "r",
			Description: "run fast, break things...",
			Action:      cmd.Run,
		},
		{
			Name:        "child",
			ShortName:   "c",
			Description: "children... how annoying!",
			Action:      cmd.Child,
		},
	}

	return app

}
