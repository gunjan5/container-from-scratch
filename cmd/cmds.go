package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/codegangsta/cli"
)

func Run(ctx *cli.Context) error {
	_ = "breakpoint"
	command := exec.Command("/proc/self/exe", append([]string{"newroot"}, ctx.Args()[0:]...)...)

	command.SysProcAttr = &syscall.SysProcAttr{ //add some namespaces: UTS, PID, MNT
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
	return nil
}

func Child(ctx *cli.Context) error {
	check(syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))
	check(os.MkdirAll("rootfs/oldrootfs", 0700))
	check(syscall.PivotRoot("rootfs", "rootfs/oldrootfs"))
	check(os.Chdir("/"))

	command := exec.Command(ctx.Args()[0], ctx.Args()[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
	return nil
}

func NewRoot(ctx *cli.Context) error {

	fmt.Println(ctx.Args()[:])

	check(os.Chdir("./OSimages/" + ctx.Args()[0]))

	if err := syscall.Chroot("."); err != nil {
		fmt.Errorf("ERROR: Chroot error ", err)
		os.Exit(1)
	}

	command := exec.Command(ctx.Args()[1], ctx.Args()[2:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
	return nil

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
