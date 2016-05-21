package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	cmd "github.com/gunjan5/container-from-scratch/cmd"
)

func main() {

	app := makeCmd()
	app.Run(os.Args)

	// switch os.Args[1] {
	// case "run":
	// 	parent()
	// case "child":
	// 	child()
	// default:
	// 	panic("wat should I do")
	//}
}

func makeCmd() *cli.App {

	app := cli.NewApp()
	app.Name = "greet"
	app.Usage = "meowwwwwwww out loud"
	// app.Action = func(c *cli.Context) error {
	// 	fmt.Println("Meow at the world", c.Args()[0])
	// 	return nil
	// }

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

// func parent() {
// 	command := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
// 	command.SysProcAttr = &syscall.SysProcAttr{
// 		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
// 	}
// 	command.Stdin = os.Stdin
// 	command.Stdout = os.Stdout
// 	command.Stderr = os.Stderr

// 	if err := command.Run(); err != nil {
// 		fmt.Println("ERROR", err)
// 		os.Exit(1)
// 	}
// }

// func child() {
// 	must(syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))
// 	must(os.MkdirAll("rootfs/oldrootfs", 0700))
// 	must(syscall.PivotRoot("rootfs", "rootfs/oldrootfs"))
// 	must(os.Chdir("/"))

// 	command := exec.Command(os.Args[2], os.Args[3:]...)
// 	command.Stdin = os.Stdin
// 	command.Stdout = os.Stdout
// 	command.Stderr = os.Stderr

// 	if err := command.Run(); err != nil {
// 		fmt.Println("ERROR", err)
// 		os.Exit(1)
// 	}
// }

// func must(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
