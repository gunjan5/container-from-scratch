package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func Run(args []string) error {
	_ = "breakpoint"
	command := exec.Command("/proc/self/exe", append([]string{"newroot"}, args[0:]...)...)

	command.SysProcAttr = &syscall.SysProcAttr{ //add some namespaces: UTS, PID, MNT
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return fmt.Errorf("ERROR executing the command calling NewRoot: ", err)
	}

	return nil

}

func NewRoot(args []string) error {

	fmt.Println(args)

	check(os.Chdir("./OSimages/" + args[0]))

	if err := syscall.Chroot("."); err != nil {
		return fmt.Errorf("ERROR: Chroot error ", err)
	}

	command := exec.Command(args[1], args[2:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return fmt.Errorf("ERROR while running the command inside the container with chroot: ", err)
	}
	return nil

}

func Child(args []string) error {
	check(syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))
	check(os.MkdirAll("rootfs/oldrootfs", 0700))
	check(syscall.PivotRoot("rootfs", "rootfs/oldrootfs"))
	check(os.Chdir("/"))

	command := exec.Command(args[0], args[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return fmt.Errorf("ERROR while running the command inside the container with pivot_root: ", err)
	}
	return nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
