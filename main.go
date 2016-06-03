package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/gunjan5/container-from-scratch/cmd"
)

func main() {

	app := makeCmd()

	if len(os.Args) > 0 {
		app.Run(os.Args)
	}

}

func makeCmd() *cli.App {

	app := cli.NewApp()
	app.Name = "CFS"
	app.Usage = "sudo ./cfs <action_command> <OS_image> <command_to_run_inside_the_container>"
	app.Version = "0.0.2"

	fmt.Println(os.Args)

	app.Commands = []cli.Command{
		{
			Name:        "server",
			ShortName:   "s",
			Description: "Start the REST server for CFS",
			Action:      cmd.Serve,
		},
		{
			Name:        "run",
			ShortName:   "r",
			Description: "run a container with a task",
			Action:      cmd.Run,
		},
		{
			Name:        "newroot",
			ShortName:   "n",
			Description: "Chroot. Not meant for direct usage",
			Action:      cmd.NewRoot,
		},
		{
			Name:        "child",
			ShortName:   "c",
			Description: "child process called by run, not meant for direct usage",
			Action:      cmd.Child,
		},
	}

	return app

}
