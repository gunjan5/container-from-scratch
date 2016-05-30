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
	app.Name = "CFS"
	app.Usage = "sudo ./cfs <action_command> <OS_image> <command_to_run_inside_the_container>"
	app.Version = "0.0.2"

	fmt.Println(os.Args)

	app.Commands = []cli.Command{
		{
			Name:        "run",
			ShortName:   "r",
			Description: "run a container with task",
			Action:      cmd.Run,
		},
		{
			Name:        "child",
			ShortName:   "c",
			Description: "child process called by run, not meant for direct usage",
			Action:      cmd.Child,
		},
		{
			Name:        "newroot",
			ShortName:   "n",
			Description: "Chroot. Not meant for direct usage",
			Action:      cmd.NewRoot,
		},
	}

	return app

}
