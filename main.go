// TO DO:
// Play around with syscall.CLONE_CGROUP and syscall.CLONE_NEWUSER (permissions issues? Try man cgroup_namespaces)

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		panic("Please use the run or child argument followed by the command to run.")
	}
}

func parent() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	// This holds the cloneflags which  will be passed to the clone syscall that gets called during cmd.Run()
	// The clone syscall creates a new process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// man clone - for more clone flags.
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())
}

func child() {
	fmt.Printf("running %v as PID %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// In order for this to work you'll need to create the rootfs folder and copy a file system there. Docker pull ubuntu
	// then cp -R /var/lib/docker/overlay2/<layer> /home/rootfs should do the trick
    must(syscall.Chroot("/home/rootfs"))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}